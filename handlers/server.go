package handlers

import (
	"log"

	"github.com/go-openapi/loads"
	"github.com/sul-dlss-labs/taco/authorization"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/validators"
)

// BuildAPI create new service API
func BuildAPI(database db.Database, storage storage.Storage, identifierService identifier.Service, depositValidator validators.ResourceValidator, updateValidator validators.ResourceValidator, fileValidator validators.ResourceValidator) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RemoteUserAuth = func(identifier string) (*authorization.Agent, error) {
		return &authorization.Agent{Identifier: identifier}, nil
	}
	api.RetrieveResourceHandler = NewRetrieveResource(database)
	api.DeleteResourceHandler = NewDeleteResource(database)
	api.RetrieveFileHandler = NewRetrieveFile(database, storage)
	api.DepositResourceHandler = NewDepositResource(database, depositValidator, identifierService)
	api.UpdateResourceHandler = NewUpdateResource(database, updateValidator)
	api.DepositFileHandler = NewDepositFile(database, storage, fileValidator, identifierService)
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
