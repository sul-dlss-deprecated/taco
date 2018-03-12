package handlers

import (
	"encoding/json"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(database db.Database) operations.UpdateResourceHandler {
	return &updateResourceEntry{database: database}
}

type updateResourceEntry struct {
	database db.Database
	stream   streaming.Stream
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(params operations.UpdateResourceParams) middleware.Responder {
	validator := validators.NewUpdateResourceValidator(d.database)
	if err := validator.ValidateResource(params.Payload); err != nil {
		return operations.NewUpdateResourceUnprocessableEntity()
	}

	resource, err := d.database.Read(params.ID)

	if err == nil {
		if err = d.updateResource(resource.ID, params); err != nil {
			panic(err)
		}

		// if err = d.addToStream(&resource.ID); err != nil {
		// 	panic(err)
		// }

		response := &models.ResourceResponse{ID: params.ID}
		return operations.NewUpdateResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

func (d *updateResourceEntry) updateResource(resourceID string, params operations.UpdateResourceParams) error {
	resource := d.persistableResourceFromParams(resourceID, params)
	return d.database.Update(resource)
}

func (d *updateResourceEntry) persistableResourceFromParams(resourceID string, params operations.UpdateResourceParams) map[string]interface{} {
	return map[string]interface{}{
		"id":        resourceID,
		"attype":    params.Payload.AtType,
		"atcontext": params.Payload.AtContext,
		"access":    params.Payload.Access,
		"label":     params.Payload.Label,
		"preserve":  params.Payload.Preserve,
		"publish":   params.Payload.Publish,
		"sourceid":  params.Payload.SourceID,
	}
}

func (d *updateResourceEntry) addToStream(id *string) error {
	message, err := json.Marshal(id)
	if err != nil {
		return err
	}

	if err := d.stream.SendMessage(string(message)); err != nil {
		return err
	}
	return nil
}
