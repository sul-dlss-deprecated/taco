package handlers

import (
	"log"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"

	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(database db.Database, validator validators.ResourceValidator, identifierService identifier.Service, authService authorization.Service) operations.DepositResourceHandler {
	return &depositResource{
		database:          database,
		validator:         validator,
		identifierService: identifierService,
		authService:       authService,
	}
}

type depositResource struct {
	database          db.Database
	validator         validators.ResourceValidator
	identifierService identifier.Service
	authService       authorization.Service
}

// Handle the create resource request
func (d *depositResource) Handle(params operations.DepositResourceParams, agent *authorization.Agent) middleware.Responder {
	resource := datautils.NewResource(params.Payload.(map[string]interface{}))

	if !d.authService.CanCreateResourceOfType(agent.Identifier, resource.Type()) {
		log.Printf("Agent %s is not permitted to create a resource of type %s", agent, resource.Type())
		return operations.NewDepositResourceUnauthorized()
	}

	if errors := d.validator.ValidateResource(resource); errors != nil {
		return operations.NewDepositResourceUnprocessableEntity().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	externalID, err := d.identifierService.Mint(resource)
	if err != nil {
		panic(err)
	}

	uuid, err := identifier.NewUUIDService().Mint(resource)
	if err != nil {
		panic(err)
	}

	resource = resource.
		WithID(uuid).
		WithExternalIdentifier(externalID).
		WithVersion(1)
	(*resource.Administrative())["created"] = time.Now().UTC().Format(time.RFC3339)

	if err := d.database.Insert(resource); err != nil {
		panic(err)
	}

	response := datautils.JSONObject{"id": externalID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}
