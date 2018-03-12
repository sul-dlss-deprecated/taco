package handlers

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

// BuildAPI create new service API
func BuildAPI(database db.Database) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RetrieveResourceHandler = NewRetrieveResource(database)
	api.DepositResourceHandler = NewDepositResource(database)
	// api.UpdateResourceHandler = NewUpdateResource(rt)
	// api.DepositFileHandler = NewDepositFile(rt)
	api.HealthCheckHandler = NewHealthCheck()
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
