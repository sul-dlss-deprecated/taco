package handlers

import (
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
	newResource := datautils.NewResource(params.Payload.(map[string]interface{})).WithID(id)

	if errors := d.validator.ValidateResource(newResource); errors != nil {
		return operations.NewUpdateResourceUnprocessableEntity().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	existingResource, err := d.database.Read(id)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	newResource = datautils.NewResource(d.mergeJSON(&existingResource.JSON, &newResource.JSON)).WithID(id)

	err = d.database.Insert(newResource)
	if err != nil {
		panic(err)
	}

	response := map[string]interface{}{"id": id}
	return operations.NewUpdateResourceOK().WithPayload(response)
}

// Merges multiple JSONObjects. Overritting in order.
func (d *updateResourceEntry) mergeJSON(maps ...*datautils.JSONObject) datautils.JSONObject {
	result := make(datautils.JSONObject)
	for _, m := range maps {
		for k, v := range *m {
			switch v.(type) {
			case datautils.JSONObject:
				if _, ok := result[k]; ok {
					x := v.(datautils.JSONObject)
					result[k] = d.mergeJSON(result.GetObj(k), &x)
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
