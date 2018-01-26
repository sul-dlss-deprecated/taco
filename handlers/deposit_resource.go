package handlers

import (
	"encoding/json"

	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
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
	requestID, _ := uuid.NewRandom()

	// TODO: This should be a DRUID
	resourceID, _ := uuid.NewRandom()
	state := "deposited"
	response := &models.DepositNewResourceOKBody{RequestID: requestID.String(), State: state, ID: resourceID.String()}
	// TODO: We probably want a different struct for the kafka message. It must have the original parameters
	message, _ := json.Marshal(response)
	d.rt.Stream().SendMessage(string(message))

	return operations.NewDepositNewResourceOK().WithPayload(response)
}
