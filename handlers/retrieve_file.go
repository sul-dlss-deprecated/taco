package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/storage"
)

// NewRetrieveFile returns a pre-signed link to a requested file
func NewRetrieveFile(repository db.Database, storage storage.Storage) operations.RetrieveFileHandler {
	return &retrieveFileEntry{repository: repository, storage: storage}
}

// retrieveFileEntry handles a request for finding & returning an file
type retrieveFileEntry struct {
	repository db.Database
	storage    storage.Storage
}

// Handle the retrieve file request
func (d *retrieveFileEntry) Handle(params operations.RetrieveFileParams) middleware.Responder {
	resource, err := d.repository.RetrieveLatest(params.ID)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	// TODO: validate that this is a file type https://github.com/sul-dlss-labs/taco/issues/214
	// TODO: check permissions https://github.com/sul-dlss-labs/taco/pull/81

	signedURL, err := d.storage.CreateSignedURL(resource.JSON.GetS("file-location"))
	if err != nil {
		panic(err)
	}
	return operations.NewRetrieveFileFound().WithLocation(*signedURL)
}
