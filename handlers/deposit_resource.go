package handlers

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"

	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(database db.Database, stream streaming.Stream, validator validators.ResourceValidator, identifierService identifier.Service) operations.DepositResourceHandler {
	return &depositResource{
		database:          database,
		stream:            stream,
		validator:         validator,
		identifierService: identifierService,
	}
}

type depositResource struct {
	database          db.Database
	stream            streaming.Stream
	validator         validators.ResourceValidator
	identifierService identifier.Service
}

// Handle the delete entry request
func (d *depositResource) Handle(params operations.DepositResourceParams) middleware.Responder {
	resource := datautils.NewResource(params.Payload.(map[string]interface{}))

	if errors := d.validator.ValidateResource(resource); errors != nil {
		return operations.NewDepositResourceUnprocessableEntity().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	resourceID, err := d.identifierService.Mint()
	if err != nil {
		panic(err)
	}
	resource.JSON["id"] = resourceID
	resource.JSON["version"] = 1

	if err = d.database.Insert(resource); err != nil {
		panic(err)
	}

	if err := d.stream.SendMessage(resourceID); err != nil {
		panic(err)
	}

	response := map[string]interface{}{"id": resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositResource) addToStream(id *string) error {
	message, err := json.Marshal(id)
	if err != nil {
		return err
	}
	if d.stream == nil {
		log.Printf("Stream is nil")
	}
	return d.stream.SendMessage(string(message))
}
