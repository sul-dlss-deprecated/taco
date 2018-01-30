package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
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
func (d *depositFileEntry) Handle(params operations.DepositFileParams, agent *authorization.Agent) middleware.Responder {
	authService := authorization.NewService(agent)
	if !authService.CanCreateResourceOfType(datautils.FileType) {
		log.Printf("Agent %s is not permitted to create a resource of type %s", agent, datautils.FileType)
		return operations.NewDepositResourceUnauthorized()
	}

	validator := validators.NewDepositFileValidator(d.database)
	if err := validator.ValidateResource(params.Upload.Header); err != nil {
		return operations.NewDepositFileInternalServerError() // TODO: need a better error
	}

	externalID, err := d.identifierService.Mint()
	if err != nil {
		panic(err)
	}

	uuid, err := identifier.NewUUIDService().Mint()
	if err != nil {
		panic(err)
	}

	upload := d.paramsToFile(params)
	location, err := d.copyFileToStorage(externalID, upload)
	if err != nil {
		panic(err)
	}

	resource := d.buildPersistableResource(upload.Metadata, location)
	resource = resource.
		WithExternalIdentifier(externalID).
		WithVersion(1).
		WithID(uuid)

	if err := d.database.Insert(resource); err != nil {
		panic(err)
	}

	url := &operations.RetrieveResourceURL{ID: externalID}
	response := datautils.JSONObject{"id": externalID}

	return operations.NewDepositResourceCreated().
		WithLocation(url.String()).
		WithPayload(response)
}

func (d *depositFileEntry) paramsToFile(params operations.DepositFileParams) *datautils.File {
	file := params.Upload
	fileHeader := file.Header
	metadata := datautils.FileMetadata{
		Filename:    fileHeader.Filename,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	return datautils.NewFile(metadata, file.Data)
}

func (d *depositFileEntry) copyFileToStorage(id string, file *datautils.File) (*string, error) {
	log.Printf("Saving file \"%s\" with content-type: %s to: %s",
		file.Metadata.Filename,
		file.Metadata.ContentType,
		id)
	return d.storage.UploadFile(id, file)
}

func (d *depositFileEntry) buildPersistableResource(metadata datautils.FileMetadata, location *string) *datautils.Resource {
	resource := NewFile()
	identification := resource.Identification()
	(*identification)["filename"] = metadata.Filename
	return resource.WithMimeType(metadata.ContentType).
		WithLabel(metadata.Filename).
		WithFileLocation(*location)
}
