package handlers

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(rt *taco.Runtime) operations.UpdateResourceHandler {
	return &updateResourceEntry{rt: rt}
}

type updateResourceEntry struct {
	rt *taco.Runtime
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(params operations.UpdateResourceParams) middleware.Responder {
	validator := validators.NewUpdateResourceValidator(d.rt.Repository())
	if err := validator.ValidateResource(params.Payload); err != nil {
		return operations.NewUpdateResourceUnprocessableEntity()
	}

	resource, err := d.rt.Repository().GetByID(params.ID)

	if err == nil {
		if err = d.updateResource(resource.ID, params); err != nil {
			panic(err)
		}

		if err = d.addToStream(&resource.ID); err != nil {
			panic(err)
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
	resource.Access = *params.Payload.Access
	resource.AtContext = *params.Payload.AtContext
	resource.AtType = *params.Payload.AtType
	resource.Label = *params.Payload.Label
	resource.Preserve = *params.Payload.Preserve
	resource.Publish = *params.Payload.Publish
	resource.SourceID = params.Payload.SourceID
	return resource
}

func (d *updateResourceEntry) addToStream(id *string) error {
	message, err := json.Marshal(id)
	if err != nil {
		return err
	}
	if d.rt.Stream() == nil {
		log.Printf("Stream is nil")
	}
	if err := d.rt.Stream().SendMessage(string(message)); err != nil {
		return err
	}
	return nil
}
