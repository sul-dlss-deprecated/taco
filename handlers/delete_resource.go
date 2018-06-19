package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/storage"
)

// NewDeleteResource -- Accepts requests to remove a resource and pushes them to Kinesis.
func NewDeleteResource(repository db.Database, storage storage.Storage) operations.DeleteResourceHandler {
	return &deleteResourceEntry{
		repository: repository,
		storage:    storage,
	}
}

type deleteResourceEntry struct {
	repository db.Database
	storage    storage.Storage
}

// Handle the delete entry request. If the resource being deleted is a file,
// Also delete the associated binary from s3
func (d *deleteResourceEntry) Handle(params operations.DeleteResourceParams) middleware.Responder {
	var err error
	if params.Version != nil {
		resource, err := d.repository.RetrieveVersion(params.ID, params.Version)
		if err == nil {
			d.deleteVersion(resource)
		}
	} else {
		err = d.deleteAllVersions(params.ID)
	}

	if err != nil {
		if _, ok := err.(*db.RecordNotFound); ok {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	return operations.NewDeleteResourceNoContent()
}

func (d *deleteResourceEntry) deleteVersion(resource *datautils.Resource) error {
	if resource.IsFile() {
		if err := d.storage.RemoveFile(resource.FileLocation()); err != nil {
			return err
		}
	}

	if err := d.repository.DeleteVersion(resource.ID()); err != nil {
		return err
	}
	return nil
}

// deleteAllVersions removes all versions with the given external id
func (d *deleteResourceEntry) deleteAllVersions(externalID string) error {
	resource, err := d.repository.RetrieveLatest(externalID)
	if err != nil {
		return err
	}
	// Delete all versions of the resource
	for resource != nil {
		err = d.deleteVersion(resource)
		if err != nil {
			return err
		}

		// retrieve the next resource
		resource, err = d.repository.RetrieveLatest(externalID)
		if err != nil {
			if _, ok := err.(*db.RecordNotFound); !ok {
				panic(err)
			}
		}

	}
	return nil
}
