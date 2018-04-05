package handlers

import (
	"log"
	"mime/multipart"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositFile -- Accepts requests to create a file and pushes it to s3.
func NewDepositFile(database db.Database, uploader storage.Storage, validator validators.ResourceValidator, identifierService identifier.Service, authService authorization.Service) operations.DepositFileHandler {
	return &depositFileEntry{database: database,
		storage:           uploader,
		validator:         validator,
		identifierService: identifierService,
		authService:       authService,
	}
}

type depositFileEntry struct {
	database          db.Database
	storage           storage.Storage
	identifierService identifier.Service
	validator         validators.ResourceValidator
	authService       authorization.Service
}

// Handle the deposit file request
func (d *depositFileEntry) Handle(params operations.DepositFileParams, agent *authorization.Agent) middleware.Responder {
	if !d.authService.CanCreateResourceOfType(agent, datautils.FileType) {
		log.Printf("Agent %s is not permitted to create a resource of type %s", agent, datautils.FileType)
		return operations.NewDepositResourceUnauthorized()
	}

	upload := d.buildFile(params.Upload.Header, params.Upload.Data)
	resource := d.buildPersistableResource(params.FilesetID, upload.Metadata)

	if errors := d.validator.ValidateResource(resource); errors != nil {
		return operations.NewDepositFileNotFound().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	externalID, err := d.identifierService.Mint(resource)
	if err != nil {
		panic(err)
	}

	uuid, err := identifier.NewUUIDService().Mint(resource)
	if err != nil {
		panic(err)
	}

	location, err := d.copyFileToStorage(externalID, upload)
	if err != nil {
		panic(err)
	}

	resource = resource.
		WithExternalIdentifier(externalID).
		WithVersion(1).
		WithFileLocation(*location).
		WithID(uuid)

	if err := d.database.Insert(resource); err != nil {
		panic(err)
	}

	// TODO: return file location: https://github.com/sul-dlss-labs/taco/issues/160
	response := datautils.JSONObject{"id": externalID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositFileEntry) buildFile(fileHeader *multipart.FileHeader, data multipart.File) *datautils.File {
	metadata := datautils.FileMetadata{
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	return datautils.NewFile(metadata, data)
}

func (d *depositFileEntry) copyFileToStorage(id string, file *datautils.File) (*string, error) {
	log.Printf("Saving file \"%s\" with content-type: %s to: %s",
		file.Metadata.Filename,
		file.Metadata.ContentType,
		id)
	return d.storage.UploadFile(id, file)
}

func (d *depositFileEntry) buildPersistableResource(FilesetID string, metadata datautils.FileMetadata) *datautils.Resource {
	resource := NewFile()
	identification := resource.Identification()
	(*identification)["filename"] = metadata.Filename
	structural := resource.Structural()
	(*structural)["isContainedBy"] = FilesetID
	return resource.WithMimeType(metadata.ContentType).
		WithLabel(metadata.Filename)
}
