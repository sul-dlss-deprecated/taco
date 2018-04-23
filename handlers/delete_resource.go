package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// NewDeleteResource -- Accepts requests to remove a resource and pushes them to Kinesis.
func NewDeleteResource(repository db.Database) operations.DeleteResourceHandler {
	return &deleteResourceEntry{
		repository: repository,
	}
}

type deleteResourceEntry struct {
	repository db.Database
	// s3
}

// Handle the delete entry request
// TODO: Delete from S3
func (d *deleteResourceEntry) Handle(params operations.DeleteResourceParams) middleware.Responder {
	if err := d.repository.DeleteAllVersions(params.ID); err != nil {
		panic(err)
	}

	return operations.NewDeleteResourceNoContent()
}
