package handlers

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// BuildAPI create new service API
func BuildAPI(rt *taco.Runtime) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RetrieveResourceHandler = NewRetrieveResource(rt)
	api.DepositNewResourceHandler = NewDepositResource(rt)
	api.HealthCheckHandler = NewHealthCheck(rt)
	return api
}

func swaggerSpec() *loads.Document {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
