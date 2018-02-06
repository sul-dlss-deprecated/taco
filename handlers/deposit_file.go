package handlers

import (
	"log"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
)

const atContext = "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld"
const fileType = "http://sdr.sul.stanford.edu/contexts/sdr3-file.jsonld"

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(rt *taco.Runtime) operations.DepositFileHandler {
	return &depositFileEntry{rt: rt}
}

type depositFileEntry struct {
	rt *taco.Runtime
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(params operations.DepositFileParams) middleware.Responder {
	id, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}

	location, err := d.copyFileToStorage(id, params.Upload)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return operations.NewDepositFileInternalServerError()
	}

	log.Printf("The location of the file is: %s", *location)

	if err := d.createFileResource(id, params.Upload.Header.Filename); err != nil {
		log.Printf("[ERROR] %s", err)
		return operations.NewDepositFileInternalServerError()
	}
	// TODO: return file location: https://github.com/sul-dlss-labs/taco/issues/160
	return operations.NewDepositResourceCreated().WithPayload(&models.ResourceResponse{ID: id})
}

func (d *depositFileEntry) copyFileToStorage(id string, file runtime.File) (*string, error) {
	return d.rt.FileStorage().UploadFile(id, file)
}

func (d *depositFileEntry) createFileResource(resourceID string, filename string) error {
	resource := d.buildPersistableResource(resourceID, filename)
	return d.rt.Repository().SaveItem(resource)
}

func (d *depositFileEntry) buildPersistableResource(resourceID string, filename string) *persistence.Resource {
	resource := &persistence.Resource{ID: resourceID}
	resource.Access = "private"
	resource.AtContext = atContext
	resource.AtType = fileType
	resource.Label = filename
	resource.Preserve = false
	resource.Publish = false
	return resource
}
