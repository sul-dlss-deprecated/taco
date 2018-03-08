package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(database db.Database) operations.RetrieveResourceHandler {
	return &retrieveResource{database: database}
}

// resourceEntry handles a request for finding & returning an entry
type retrieveResource struct {
	database db.Database
}

// Handle the delete entry request
func (d *retrieveResource) Handle(params operations.RetrieveResourceParams) middleware.Responder {
	item, err := d.database.Read(params.ID)
	if err == nil {
		// TODO: expand this mapping
		response := buildResponse(item)
		return operations.NewRetrieveResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

// TODO: expand this mapping
func buildResponse(resource *persistence.Resource) models.Resource {
	return resource
}
