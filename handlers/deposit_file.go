package handlers

import (
	"log"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/uploaded"
	"github.com/sul-dlss-labs/taco/validators"
)

const atContext = "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld"
const fileType = "http://sdr.sul.stanford.edu/contexts/sdr3-file.jsonld"

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(database db.Database, uploader storage.Storage) operations.DepositFileHandler {
	return &depositFileEntry{database: database, storage: uploader}
}

type depositFileEntry struct {
	database db.Database
	storage  storage.Storage
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(params operations.DepositFileParams) middleware.Responder {
	validator := validators.NewDepositFileValidator(d.database)
	if err := validator.ValidateResource(params.Upload.Header); err != nil {
		return operations.NewDepositFileInternalServerError() // TODO: need a better error
	}

	id, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}

	location, err := d.copyFileToStorage(id, params.Upload)
	if err != nil {
		panic(err)
	}

	log.Printf("The location of the file is: %s", *location)

	if err := d.createFileResource(id, params.Upload.Header.Filename); err != nil {
		panic(err)
	}
	// TODO: return file location: https://github.com/sul-dlss-labs/taco/issues/160
	response := map[string]interface{}{"id": id}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositFileEntry) copyFileToStorage(id string, file runtime.File) (*string, error) {
	filename := file.Header.Filename
	contentType := file.Header.Header.Get("Content-Type")
	log.Printf("Saving file \"%s\" with content-type: %s", filename, contentType)

	upload := uploaded.NewFile(filename, contentType, file.Data)
	return d.storage.UploadFile(id, upload)
}

func (d *depositFileEntry) createFileResource(resourceID string, filename string) error {
	resource := d.buildPersistableResource(resourceID, filename)
	return d.database.Insert(resource)
}

func (d *depositFileEntry) buildPersistableResource(resourceID string, filename string) datautils.Resource {
	// TODO: Expand here if we need to set any default properties on the file
	identification := map[string]interface{}{"filename": filename, "identifier": resourceID, "sdrUUID": resourceID}
	return datautils.Resource{"id": resourceID, "identification": identification}
}
