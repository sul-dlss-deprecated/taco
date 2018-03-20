package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(database db.Database) operations.RetrieveResourceHandler {
	return &retrieveResource{database: database}
}

// retrieveResource handles a request for finding & returning an entry
type retrieveResource struct {
	database db.Database
}

// Handle the retrieve resource request
func (d *retrieveResource) Handle(params operations.RetrieveResourceParams) middleware.Responder {
	resource, err := d.database.Read(params.ID)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	return operations.NewRetrieveResourceOK().WithPayload(resource.JSON)
}
