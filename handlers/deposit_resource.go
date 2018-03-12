package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/resource"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(database db.Database) operations.DepositResourceHandler {
	return &depositResource{database: database}
}

type depositResource struct {
	database db.Database
	// stream
	// validators
}

// Handle the delete entry request
func (d *depositResource) Handle(params operations.DepositResourceParams) middleware.Responder {
	fmt.Printf("%+v\n", params)
	/*
		validator := validators.NewDepositResourceValidator(repository)
		if err := validator.ValidateResource(params.Payload); err != nil {
			return operations.NewDepositResourceUnprocessableEntity()
		}
	*/

	resourceID, err := identifier.NewService().Mint()
	fmt.Printf("resourceID: %s\n", resourceID)
	if err != nil {
		panic(err)
	}

	err = d.database.Insert(resource.LoadParams(resourceID, params))
	if err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	/*
		if err := d.addToStream(&resourceID); err != nil {
			// TODO: handle this with an error response
			panic(err)
		}
	*/
	response := &models.ResourceResponse{ID: resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}
