package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(rt *taco.Runtime) operations.RetrieveResourceHandler {
	return &resourceEntry{repository: rt.Repository()}
}

// resourceEntry handles a request for finding & returning an entry
type resourceEntry struct {
	repository persistence.Repository
}

// Handle the delete entry request
func (d *resourceEntry) Handle(params operations.RetrieveResourceParams) middleware.Responder {
	resource, err := d.repository.GetByID(params.ID)
	if err == nil {
		// TODO: expand this mapping
		response := buildResponse(resource)
		return operations.NewRetrieveResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

// TODO: expand this mapping
func buildResponse(resource *persistence.Resource) *models.Resource {
	return &models.Resource{
		ID:        resource.ID,
		Label:     &resource.Label,
		AtContext: &resource.AtContext,
		AtType:    &resource.AtType,
		Access:    &resource.Access,
		Preserve:  &resource.Preserve,
		Publish:   &resource.Publish}
}
