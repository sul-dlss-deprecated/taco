package handlers

import (
	"encoding/json"
	"log"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
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

// Handle the deposit request
func (d *depositResourceEntry) Handle(params operations.DepositResourceParams) middleware.Responder {
	validator := validators.NewDepositResourceValidator(d.rt.Repository())
	if err := validator.ValidateResource(params.Payload); err != nil {
		return operations.NewDepositResourceUnprocessableEntity()
	}

	sdrUUID, resourceID := d.identifierForType(params.Payload.AtType)

	if err := d.persistResource(sdrUUID, resourceID, params); err != nil {
		panic(err)
	}

	if err := d.addToStream(&resourceID); err != nil {
		panic(err)
	}

	response := &models.ResourceResponse{ID: resourceID}
	return operations.NewDepositResourceCreated().WithPayload(response)
}

// If configured to talk to the identifier service (SURI) and if the resource is
// a Collection or DRO, return the DRUID from the remote service. Otherwise
// return a uuid
func (d *depositResourceEntry) identifierForType(resourceType *strfmt.URI) (string, string) {
	sdrUUID, resourceID, err := d.rt.Identifier().Mint(*resourceType)
	if err != nil {
		panic(err)
	}
	return sdrUUID, resourceID
}

func (d *depositResourceEntry) persistResource(sdrUUID string, resourceID string, params operations.DepositResourceParams) error {
	resource := d.persistableResourceFromParams(sdrUUID, resourceID, params)
	return d.rt.Repository().CreateItem(resource)
}

func (d *depositResourceEntry) persistableResourceFromParams(sdrUUID string, resourceID string, params operations.DepositResourceParams) *persistence.Resource {
	resource := &persistence.Resource{SdrUUID: sdrUUID, Identifier: resourceID}
	log.Printf("Depositing %s, %s", resource.SdrUUID, resourceID)
	resource.Access = *params.Payload.Access
	resource.AtContext = *params.Payload.AtContext
	resource.AtType = *params.Payload.AtType
	resource.Label = *params.Payload.Label
	resource.Preserve = *params.Payload.Preserve
	resource.Publish = *params.Payload.Publish
	resource.SourceID = params.Payload.SourceID
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
