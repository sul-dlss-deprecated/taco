package handlers

import (
	"encoding/json"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
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
	json, err := json.Marshal(params.Payload)
	if err != nil {
		panic(err)
	}
	if err := d.validator.ValidateResource(string(json[:])); err != nil {
		return operations.NewUpdateResourceUnprocessableEntity()
	}

	id := params.ID
	newResource := datautils.NewResource(params.Payload.(map[string]interface{}))

	existingResource, err := d.database.Read(id)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	// newResource.MergeResource(existingResource)
	newResource = d.mergeResources(*existingResource, newResource)
	newResource["id"] = id
	if err != nil {
		panic(err)
	}

	err = d.database.Insert(newResource)
	if err != nil {
		panic(err)
	}

	response := map[string]interface{}{"id": id}
	return operations.NewUpdateResourceOK().WithPayload(response)
}

// Merges multiple map[string]interface{} objects. Overritting in order.
func (d *updateResourceEntry) mergeResources(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			switch v.(type) {
			case map[string]interface{}:
				if _, ok := result[k]; ok {
					result[k] = d.mergeResources(result[k].(map[string]interface{}), v.(map[string]interface{}))
				} else {
					result[k] = v
				}
			default:
				result[k] = v
			}
		}
	}
	return result
}
