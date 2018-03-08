package handlers

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/justinas/alice"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/middleware"
)

// BuildHandler sets up the middleware that wraps the API
func Setup(api *operations.TacoAPI) http.Handler {
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
