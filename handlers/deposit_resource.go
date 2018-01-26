package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/uuid"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// NewDepositResource -- Accepts requests to create resource and pushes them to Kinesis.
func NewDepositResource(rt *taco.Runtime) operations.DepositNewResourceHandler {
	return &depositResourceEntry{}
}

type depositResourceEntry struct{}

// Handle the delete entry request
func (d *depositResourceEntry) Handle(params operations.DepositNewResourceParams) middleware.Responder {
	requestID, _ := uuid.NewRandom()

	// TODO: This should be a DRUID
	resourceID, _ := uuid.NewRandom()

	response := &models.Resource{RequestID: requestID.String(), ID: resourceID.String()}
	return operations.NewDepositNewResourceOK().WithPayload(response)
}
