package handlers

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(database db.Database, stream streaming.Stream, validator validators.ResourceValidator) operations.UpdateResourceHandler {
	return &updateResourceEntry{database: database, stream: stream, validator: validator}
}

type updateResourceEntry struct {
	database  db.Database
	stream    streaming.Stream
	validator validators.ResourceValidator
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(params operations.UpdateResourceParams) middleware.Responder {
	id := params.ID
	resource := d.persistableResourceFromParams(id, params.Payload)

	if errors := d.validator.ValidateResource(resource); errors != nil {
		return operations.NewUpdateResourceUnprocessableEntity().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	_, err := d.database.Read(id)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	if err = d.database.Update(resource); err != nil {
		panic(err)
	}

	if err = d.addToStream(id); err != nil {
		panic(err)
	}

	response := map[string]interface{}{"id": id}
	return operations.NewUpdateResourceOK().WithPayload(response)
}

func (d *updateResourceEntry) persistableResourceFromParams(resourceID string, data models.Resource) *datautils.Resource {
	resource := datautils.NewResource(data.(map[string]interface{}))
	// This ensures they have the same id in the document as in the query param
	resource.JSON["id"] = resourceID
	return resource
}

func (d *updateResourceEntry) addToStream(id string) error {
	message, err := json.Marshal(id)
	if err != nil {
		return err
	}
	if d.stream == nil {
		log.Printf("Stream is nil")
	}
	return d.stream.SendMessage(string(message))
}
