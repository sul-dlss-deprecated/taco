package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(database db.Database, authService authorization.Service) operations.RetrieveResourceHandler {
	return &retrieveResource{database: database, authService: authService}
}

// retrieveResource handles a request for finding & returning an entry
type retrieveResource struct {
	database    db.Database
	authService authorization.Service
}

// Handle the retrieve resource request
func (d *retrieveResource) Handle(params operations.RetrieveResourceParams, agent *authorization.Agent) middleware.Responder {
	var resource *datautils.Resource
	var err error
	if params.Version != nil {
		resource, err = d.database.RetrieveVersion(params.ID, params.Version)
	} else {
		resource, err = d.database.RetrieveLatest(params.ID)
	}

	if err != nil {
		if _, ok := err.(*db.RecordNotFound); ok {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	if !d.authService.CanRetrieveResource(agent, resource) {
		log.Printf("Agent %s is not permitted to retrieve this resource %s", agent, params.ID)
		return operations.NewRetrieveResourceUnauthorized()
	}

	return operations.NewRetrieveResourceOK().WithPayload(resource.JSON)
}
