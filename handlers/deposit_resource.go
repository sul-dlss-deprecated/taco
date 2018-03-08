package handlers

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/validators"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(rt *taco.Runtime) operations.DepositResourceHandler {
	return &depositResourceEntry{rt: rt}
}

type depositResourceEntry struct {
	rt *taco.Runtime
}

// Handle the delete entry request
func (d *depositResourceEntry) Handle(params operations.DepositResourceParams) middleware.Responder {
	validator := validators.NewDepositResourceValidator(d.rt.Repository())
	if err := validator.ValidateResource(params.Payload); err != nil {
		return operations.NewDepositResourceUnprocessableEntity()
	}

	resourceID, err := identifier.NewService().Mint()
	if err != nil {
		panic(err)
	}
	if err := d.persistResource(resourceID, params); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	if err := d.addToStream(&resourceID); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	response := &models.ResourceResponse{ID: resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

func (d *depositResourceEntry) persistResource(resourceID string, params operations.DepositResourceParams) error {
	resource := d.persistableResourceFromParams(resourceID, params)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositResourceEntry) persistableResourceFromParams(resourceID string, params operations.DepositResourceParams) *persistence.Resource {
	resource := &persistence.Resource{ID: resourceID}
	resource.Access = models.ResourceAccess{Access: params.Payload.Access.Access}
	resource.AtContext = *params.Payload.AtContext
	resource.AtType = *params.Payload.AtType
	resource.Label = *params.Payload.Label
	// TODO: ResourceIdentification has no SourceID?
	//resource.Identification = models.ResourceIdentification{SourceID: params.Payload.Identification.SourceID}
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
	return nil
}
