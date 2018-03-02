package handlers

import (
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/uploaded"
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
	filename := file.Header.Filename
	contentType := file.Header.Header.Get("Content-Type")
	log.Printf("Saving file \"%s\" with content-type: %s", filename, contentType)

	upload := uploaded.NewFile(filename, contentType, file.Data)
	return d.rt.FileStorage().UploadFile(id, upload)
}

func (d *depositFileEntry) createFileResource(resourceID string, filename string) error {
	resource := d.buildPersistableResource(resourceID, filename)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositFileEntry) buildPersistableResource(resourceID string, filename string) map[string]*dynamodb.AttributeValue {
	row := map[string]*dynamodb.AttributeValue{}
	ctx := atContext
	typ := fileType
	row["atContext"] = &dynamodb.AttributeValue{S: &ctx}
	row["atType"] = &dynamodb.AttributeValue{S: &typ}
	row["Label"] = &dynamodb.AttributeValue{S: &filename}
	row[persistence.PrimaryKey] = &dynamodb.AttributeValue{S: &resourceID}

	// TODO: Do we need any of these?
	// resource.Access = "private"
	// resource.Preserve = false
	// resource.Publish = false

	return row
}
