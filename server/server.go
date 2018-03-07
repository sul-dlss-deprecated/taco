package server

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/handlers"
	"github.com/sul-dlss-labs/taco/middleware"
)

// BuildAPI create new service API
func BuildAPI() *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RetrieveResourceHandler = handlers.NewRetrieveResource()
	api.DepositResourceHandler = handlers.NewDepositResource()
	api.UpdateResourceHandler = handlers.NewUpdateResource()
	api.DepositFileHandler = handlers.NewDepositFile()
	api.HealthCheckHandler = handlers.NewHealthCheck()
	return api
}

// BuildHandler sets up the middleware that wraps the API
func BuildHandler(api *operations.TacoAPI) http.Handler {
	return alice.New(
		middleware.NewHoneyBadgerMW(),
		middleware.NewRecoveryMW(),
		middleware.NewRequestLoggerMW(),
	).Then(api.Serve(nil))
}

func swaggerSpec() *loads.Document {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
