package handlers

import (
	"encoding/json"
	"log"
	"path"
	"runtime"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
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

	if err = d.compareAndUpdateResource(id, newResource, existingResource); err != nil {
		panic(err)
	}

	if err = d.addToStream(id); err != nil {
		panic(err)
	}

	response := map[string]interface{}{"id": id}
	return operations.NewUpdateResourceOK().WithPayload(response)
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

func (d *updateResourceEntry) compareAndUpdateResource(id string, newResource datautils.Resource, existingResource *datautils.Resource) error {

	var err error

	for k, v := range newResource {
		if existingResource.HasKey(k) {
			switch v.(type) {
			case string:
				if v != existingResource.GetS(k) {
					err = d.database.UpdateString(id, k, v.(string))
				}
			case json.Number:
				if v.(json.Number).String() != existingResource.GetS(k) {
					err = d.database.UpdateString(id, k, v.(json.Number).String())
				}
			case bool:
				if v != existingResource.GetB(k) {
					err = d.database.UpdateBool(id, k, v.(bool))
				}
			case map[string]interface{}:
				// TODO: Update on nested values
				// updates := make(map[string]interface{})
				// updates[k] = d.compareMap(v, (*oldRes)[k])
			}
		} else {
			switch v.(type) {
			case string:
				err = d.database.UpdateString(id, k, v.(string))
			case json.Number:
				err = d.database.UpdateString(id, k, v.(json.Number).String())
			case bool:
				err = d.database.UpdateBool(id, k, v.(bool))
			}
		}
	}

	return err
}
