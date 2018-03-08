package handlers

import (
	"encoding/json"
	"log"
	"path"
	"runtime"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(database db.Database, stream streaming.Stream) operations.UpdateResourceHandler {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps/Resource.json")
	validator := validators.NewUpdateResourceValidator(database, schemaPath)

	return &updateResourceEntry{database: database, stream: stream, validator: validator}
}

type updateResourceEntry struct {
	database  db.Database
	stream    streaming.Stream
	validator *validators.UpdateResourceValidator
}

// Handle the update resource request
func (d *updateResourceEntry) Handle(params operations.UpdateResourceParams) middleware.Responder {
	if err := d.validator.ValidateResource(params.Payload); err != nil {
		return operations.NewUpdateResourceUnprocessableEntity()
	}

	resource, err := d.database.Read(params.ID)

	if err == nil {
		if err = d.updateResource(resource.ID(), params); err != nil {
			panic(err)
		}

		if err = d.addToStream(id); err != nil {
			panic(err)
		}

		response := map[string]interface{}{"id": id}
		return operations.NewUpdateResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

func (d *updateResourceEntry) updateResource(resourceID string, params operations.UpdateResourceParams) error {
	resource := d.persistableResourceFromParams(resourceID, params)
	return d.database.UpdateItem(resource)
}

func (d *updateResourceEntry) persistableResourceFromParams(resourceID string, params operations.UpdateResourceParams) persistence.Resource {
	resource := persistence.Resource{"id": resourceID}
	// resource.Access = *params.Payload.Access
	// resource.AtContext = *params.Payload.AtContext
	// resource.AtType = *params.Payload.AtType
	// resource.Label = *params.Payload.Label
	// resource.Preserve = *params.Payload.Preserve
	// resource.Publish = *params.Payload.Publish
	// resource.SourceID = params.Payload.SourceID
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
	if err := d.stream.SendMessage(string(message)); err != nil {
		return err
	}
}
