package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// Retrieve Resource will query DynamoDB with ID for Resource JSON
// Currently, just returns fake object for any ID.
func NewRetrieveResource() operations.RetrieveResourceHandler {
	return &resourceEntry{}
}

// resourceEntry handles a request for finding & returning an entry
type resourceEntry struct {
}

// Handle the delete entry request
func (d *resourceEntry) Handle(params operations.RetrieveResourceParams) middleware.Responder {
	id := swag.StringValue(&params.ID)
	test := new(models.Resource)
	test.ID = id
	return operations.NewRetrieveResourceOK().WithPayload(test)
}
