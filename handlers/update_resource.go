package handlers

import (
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewUpdateResource(rt *taco.Runtime) operations.UpdateResourceHandler {
	return &updateResourceEntry{rt: rt, repository: rt.Repository()}
}

type updateResourceEntry struct {
	rt         *taco.Runtime
	repository persistence.Repository
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(params operations.UpdateResourceParams) middleware.Responder {
	resource, err := d.repository.GetByID(params.ID)

	if err == nil {
		if err := d.updateResource(resource.ID, params); err != nil {
			// TODO: handle this with an error response
			panic(err)
		} else {
			log.Printf("Resource Update Successful")
		}

		response := &models.ResourceResponse{ID: params.ID}
		return operations.NewUpdateResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

func (d *updateResourceEntry) updateResource(resourceID string, params operations.UpdateResourceParams) error {
	resource := d.persistableResourceFromParams(resourceID, params)
	return d.rt.Repository().UpdateItem(resource)
}

func (d *updateResourceEntry) persistableResourceFromParams(resourceID string, params operations.UpdateResourceParams) *persistence.Resource {
	resource := &persistence.Resource{ID: resourceID}
	resource.Access = *params.Body.Access
	resource.AtContext = *params.Body.AtContext
	resource.AtType = *params.Body.AtType
	resource.Label = *params.Body.Label
	resource.Preserve = *params.Body.Preserve
	resource.Publish = *params.Body.Publish
	return resource
}
