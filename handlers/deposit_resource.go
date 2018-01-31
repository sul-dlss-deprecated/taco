package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/persistence"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(rt *taco.Runtime) operations.DepositNewResourceHandler {
	return &depositResourceEntry{rt: rt}
}

type depositResourceEntry struct {
	rt *taco.Runtime
}

// Handle the delete entry request
func (d *depositResourceEntry) Handle(params operations.DepositNewResourceParams) middleware.Responder {

	resourceID := mintID()

	if err := d.persistResource(resourceID, params); err != nil {
		// TODO: handle this with an error response
		panic(err)
	}

	response := &models.ResourceResponse{ID: resourceID}
	return operations.NewDepositNewResourceCreated().WithPayload(response)
}

func (d *depositResourceEntry) persistResource(resourceID string, params operations.DepositNewResourceParams) error {
	resource := d.persistableResourceFromParams(resourceID, params)
	fmt.Println("Saving")
	return d.rt.Repository().SaveItem(resource)
}

func (d *depositResourceEntry) persistableResourceFromParams(resourceID string, params operations.DepositNewResourceParams) *persistence.Resource {
	resource := &persistence.Resource{ID: resourceID}
	resource.Access = *params.Payload.Access
	resource.AtContext = *params.Payload.AtContext
	resource.AtType = *params.Payload.AtType
	resource.Label = *params.Payload.Label
	resource.Preserve = *params.Payload.Preserve
	resource.Publish = *params.Payload.Publish
	// resource.SourceID = *params.Payload.SourceID
	return resource
}

func mintID() string {
	// TODO: This should be a DRUID
	resourceID, _ := uuid.NewRandom()
	return resourceID.String()
}
