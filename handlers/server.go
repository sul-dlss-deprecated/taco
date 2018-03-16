package handlers

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
	"github.com/sul-dlss-labs/taco/validators"
)

// BuildAPI create new service API
func BuildAPI(database db.Database, stream streaming.Stream, storage storage.Storage, schemaDir string) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RetrieveResourceHandler = NewRetrieveResource(database)
	depositValidator := validators.NewDepositResourceValidator(database, schemaDir)
	api.DepositResourceHandler = NewDepositResource(database, stream, depositValidator)
	updateValidator := validators.NewUpdateResourceValidator(database, schemaDir)
	api.UpdateResourceHandler = NewUpdateResource(database, stream, updateValidator)
	api.DepositFileHandler = NewDepositFile(database, storage)
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
