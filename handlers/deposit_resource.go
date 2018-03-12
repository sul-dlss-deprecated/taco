package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/resource"
	"github.com/sul-dlss-labs/taco/streaming"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(database db.Database, stream streaming.Stream) operations.DepositResourceHandler {
	return &depositResource{database: database, stream: stream}
}

type depositResource struct {
	database db.Database
	stream   streaming.Stream
	// validators
}

// Handle the delete entry request
func (d *depositResource) Handle(params operations.DepositResourceParams) middleware.Responder {
	/*
		validator := validators.NewDepositResourceValidator(repository)
		if err := validator.ValidateResource(params.Payload); err != nil {
			return operations.NewDepositResourceUnprocessableEntity()
		}
	*/

	resourceID, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}

	err = d.database.Insert(resource.LoadParams(resourceID, params))
	if err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	if err := d.stream.SendMessage(resourceID); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	response := &models.ResourceResponse{ID: resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}
