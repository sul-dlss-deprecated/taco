package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// NewHealthCheck will return the service health
func NewHealthCheck() operations.HealthCheckHandler {
	return &healthCheck{}
}

type healthCheck struct{}

// Handle the health check request
func (d *healthCheck) Handle(params operations.HealthCheckParams) middleware.Responder {
	return operations.NewHealthCheckOK().WithPayload(&models.HealthCheckResponse{Status: "OK"})
}
