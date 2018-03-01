package handlers

import (
	"log"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/uploaded"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(rt *taco.Runtime) operations.DepositFileHandler {
	return &depositFileEntry{rt: rt}
}

type depositFileEntry struct {
	rt *taco.Runtime
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(params operations.DepositFileParams) middleware.Responder {
	validator := validators.NewDepositFileValidator(d.rt.Repository())
	if err := validator.ValidateResource(params.Upload.Header); err != nil {
		return operations.NewDepositFileInternalServerError() // TODO: need a better error
	}

	sdrUUID, resourceID, err := d.rt.Identifier().Mint(persistence.FileType)
	if err != nil {
		panic(err)
	}

	location, err := d.copyFileToStorage(resourceID, params.Upload)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return operations.NewDepositFileInternalServerError()
	}

	log.Printf("The location of the file is: %s", *location)

	if err := d.createFileResource(sdrUUID, resourceID, params.Upload.Header.Filename); err != nil {
		log.Printf("[ERROR] %s", err)
		return operations.NewDepositFileInternalServerError()
	}
	// TODO: return file location: https://github.com/sul-dlss-labs/taco/issues/160
	return operations.NewDepositResourceCreated().WithPayload(&models.ResourceResponse{ID: resourceID})
}

func (d *depositFileEntry) copyFileToStorage(id string, file runtime.File) (*string, error) {
	filename := file.Header.Filename
	contentType := file.Header.Header.Get("Content-Type")
	log.Printf("Saving file \"%s\" with content-type: %s", filename, contentType)

	upload := uploaded.NewFile(filename, contentType, file.Data)
	return d.rt.FileStorage().UploadFile(id, upload)
}

func (d *depositFileEntry) createFileResource(sdrUUID string, resourceID string, filename string) error {
	resource := d.buildPersistableResource(sdrUUID, resourceID, filename)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositFileEntry) buildPersistableResource(sdrUUID string, resourceID string, filename string) *persistence.Resource {
	resource := &persistence.Resource{SdrUUID: sdrUUID, Identifier: resourceID}
	// TODO: Where should Access come from/default to?
	resource.Access = "private"
	resource.AtContext = persistence.AtContext
	resource.AtType = persistence.FileType
	resource.Label = filename
	resource.Preserve = false
	resource.Publish = false
	return resource
}
