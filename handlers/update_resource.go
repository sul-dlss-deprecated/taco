package handlers

import (
	"encoding/json"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewUpdateResource -- Accepts requests to update a resource.
func NewUpdateResource(database db.Database, validator validators.ResourceValidator) operations.UpdateResourceHandler {
	return &updateResourceEntry{database: database, validator: validator}
}

type updateResourceEntry struct {
	database  db.Database
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

	existingResource, err := d.database.RetrieveLatest(id)
	if err != nil {
		if err.Error() == "not found" {
			return operations.NewRetrieveResourceNotFound()
		}
		panic(err)
	}

	merged := d.mergeJSON(&existingResource.JSON, &newResource.JSON)
	newResource = datautils.NewResource(merged).
		WithExternalIdentifier(id).   // Don't allow changing druids
		WithID(existingResource.ID()) // Ignore any passed in tacoIdentifier

	err = d.database.Insert(newResource)
	if err != nil {
		panic(err)
	}

	response := datautils.JSONObject{"id": id}
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
			case json.Number:
				// Cast "version" to int, otherswise dynamodbattribute.MarshalMap will cast it to String
				// See https://github.com/aws/aws-sdk-go-v2/issues/115
				i64, err := v.(json.Number).Int64()
				if err != nil {
					panic(err)
				}
				result[k] = int(i64)
			default:
				result[k] = v
			}
		}
	}
	return result
}
