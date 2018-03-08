package handlers

import (
	"encoding/json"
	"log"
	"path"
	"runtime"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(database db.Database, stream streaming.Stream) operations.DepositResourceHandler {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps/DepositResource.json")
	validator := validators.NewDepositResourceValidator(database, schemaPath)
	return &depositResource{database: database, stream: stream, validator: validator}
}

type depositResource struct {
	database  db.Database
	stream    streaming.Stream
	validator *validators.DepositResourceValidator
}

// Handle the delete entry request
func (d *depositResource) Handle(params operations.DepositResourceParams) middleware.Responder {
	if err := validator.ValidateResource(params.Payload); err != nil {
		return operations.NewDepositResourceUnprocessableEntity()
	}

	resourceID, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}

	err = d.database.Insert(d.loadParams(resourceID, params))
	if err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	if err := d.stream.SendMessage(resourceID); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	response := map[string]interface{}{"id": resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositResourceEntry) loadParams(resourceID string, params operations.DepositResourceParams) persistence.Resource {
	resource := persistence.Resource{"id": resourceID}
	return resource
}

func (d *depositResourceEntry) addToStream(id *string) error {
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
}
