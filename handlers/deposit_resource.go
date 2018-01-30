package handlers

import (
	"log"

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
func NewDepositResource(database db.Database, validator validators.ResourceValidator, identifierService identifier.Service) operations.DepositResourceHandler {
	return &depositResource{
		database:          database,
		validator:         validator,
		identifierService: identifierService,
	}
}

type depositResource struct {
	database          db.Database
	validator         validators.ResourceValidator
	identifierService identifier.Service
}

// Handle the create resource request
func (d *depositResource) Handle(params operations.DepositResourceParams, agent *authorization.Agent) middleware.Responder {
	resource := datautils.NewResource(params.Payload.(map[string]interface{}))

	if errors := d.validator.ValidateResource(resource); errors != nil {
		return operations.NewDepositResourceUnprocessableEntity().
			WithPayload(&models.ErrorResponse{Errors: *errors})
	}

	authService := authorization.NewService(agent)
	if !authService.CanCreateResourceOfType(resource.Type()) {
		log.Printf("Agent %s is not permitted to create a resource of type %s", agent, resource.Type())
		return operations.NewDepositResourceUnauthorized()
	}

	externalID, err := d.identifierService.Mint()
	if err != nil {
		panic(err)
	}

	uuid, err := identifier.NewUUIDService().Mint()
	if err != nil {
		panic(err)
	}

	resource = resource.
		WithID(uuid).
		WithExternalIdentifier(externalID).
		WithVersion(1)

	if err = d.database.Insert(resource); err != nil {
		panic(err)
	}

	response := datautils.JSONObject{"id": externalID}
	url := &operations.RetrieveResourceURL{ID: externalID}

	return operations.NewDepositResourceCreated().
		WithLocation(url.String()).
		WithPayload(response)
}
