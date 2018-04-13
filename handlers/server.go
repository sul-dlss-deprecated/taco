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
func BuildAPI(database db.Database, storage storage.Storage, identifierService identifier.Service) *operations.TacoAPI {
	api := operations.NewTacoAPI(swaggerSpec())
	api.RemoteUserAuth = func(identifier string) (*authorization.Agent, error) {
		return &authorization.Agent{Identifier: identifier}, nil
	}
	api.RetrieveResourceHandler = NewRetrieveResource(database)
	api.DeleteResourceHandler = NewDeleteResource(database)
	api.RetrieveFileHandler = NewRetrieveFile(database, storage)
	api.DepositResourceHandler = NewDepositResource(database, depositResourceValidator(database), identifierService)
	api.UpdateResourceHandler = NewUpdateResource(database, updateResourceValidator(database))
	api.DepositFileHandler = NewDepositFile(database, storage, depositFileValidator(database), identifierService)
	api.HealthCheckHandler = NewHealthCheck()
	return api
}

// Builds the validator for deposit file requests
func depositFileValidator(database db.Database) validators.ResourceValidator {
	return validators.NewCompositeResourceValidator(
		[]validators.ResourceValidator{
			validators.NewDepositFileValidator(database),
			structuralValidator(database),
		})
}

// Builds the validator for deposit resource requests
func depositResourceValidator(database db.Database) validators.ResourceValidator {
	return validators.NewCompositeResourceValidator(
		[]validators.ResourceValidator{
			validators.NewDepositResourceValidator(database),
			structuralValidator(database),
		})
}

// Builds the validator for update resource requests
func updateResourceValidator(database db.Database) validators.ResourceValidator {
	return validators.NewCompositeResourceValidator(
		[]validators.ResourceValidator{
			validators.NewUpdateResourceValidator(database),
			structuralValidator(database),
		})
}

// Builds the validator for structural validations.
// This is suitable for both create and update requests
func structuralValidator(database db.Database) validators.ResourceValidator {
	return validators.NewCompositeResourceValidator(
		[]validators.ResourceValidator{
			validators.NewFileStructuralValidator(database),
			validators.NewFilesetStructuralValidator(database),
			validators.NewDROStructuralValidator(database),
			validators.NewCollectionStructuralValidator(database),
			validators.NewAgreementValidator(database),
			validators.NewSequenceValidator(),
		})
}

func swaggerSpec() *loads.Document {
	// load embedded swagger file
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}
	return swaggerSpec
}
