package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/uploaded"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(database db.Database, uploader storage.Storage, identifierService identifier.Service) operations.DepositFileHandler {
	return &depositFileEntry{database: database,
		storage:           uploader,
		identifierService: identifierService}
}

type depositFileEntry struct {
	database          db.Database
	storage           storage.Storage
	identifierService identifier.Service
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(params operations.DepositFileParams) middleware.Responder {
	validator := validators.NewDepositFileValidator(d.database)
	if err := validator.ValidateResource(params.Upload.Header); err != nil {
		return operations.NewDepositFileInternalServerError() // TODO: need a better error
	}

	id, err := d.identifierService.Mint()
	if err != nil {
		panic(err)
	}

	upload := d.paramsToFile(params)
	location, err := d.copyFileToStorage(id, upload)
	if err != nil {
		panic(err)
	}

	log.Printf("The location of the file is: %s", *location)

	if err := d.createFileResource(id, upload.Metadata); err != nil {
		panic(err)
	}
	// TODO: return file location: https://github.com/sul-dlss-labs/taco/issues/160
	response := datautils.JSONObject{"id": id}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositFileEntry) paramsToFile(params operations.DepositFileParams) *uploaded.File {
	file := params.Upload
	fileHeader := file.Header
	metadata := uploaded.FileMetadata{
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	return uploaded.NewFile(metadata, file.Data)
}

func (d *depositFileEntry) copyFileToStorage(id string, file *uploaded.File) (*string, error) {
	log.Printf("Saving file \"%s\" with content-type: %s to: %s",
		file.Metadata.Filename,
		file.Metadata.ContentType,
		id)
	return d.storage.UploadFile(id, file)
}

func (d *depositFileEntry) createFileResource(resourceID string, metadata uploaded.FileMetadata) error {
	resource := d.buildPersistableResource(resourceID, metadata)
	return d.database.Insert(resource)
}

func (d *depositFileEntry) buildPersistableResource(resourceID string, metadata uploaded.FileMetadata) *datautils.Resource {
	identification := map[string]interface{}{"filename": metadata.Filename, "identifier": resourceID, "sdrUUID": resourceID}
	json := datautils.JSONObject{"id": resourceID, "identification": identification, "hasMimeType": metadata.ContentType}
	return datautils.NewResource(json)
}
