package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
	"github.com/sul-dlss-labs/taco/resource"
)

// NewRetrieveResource will query DynamoDB with ID for Resource JSON
func NewRetrieveResource(database *db.Database) operations.RetrieveResourceHandler {
	return &retrieveResource{database: database}
}

// resourceEntry handles a request for finding & returning an entry
type retrieveResource struct {
	database *db.Database
}

// Handle the delete entry request
func (d *retrieveResource) Handle(params operations.RetrieveResourceParams) middleware.Responder {
	resource, err := resource.GetByID(params.ID)
	if err == nil {
		// TODO: expand this mapping
		// response := buildResponse(resource)
		return operations.NewRetrieveResourceOK().WithPayload(response)
	} else if err.Error() == "not found" {
		return operations.NewRetrieveResourceNotFound()
	}
	panic(err)
}

// TODO: expand this mapping
func buildResponse(resource *persistence.Resource) *models.Resource {
	return &models.Resource{
		ID:        resource.ID,
		Label:     &resource.Label,
		AtContext: &resource.AtContext,
		AtType:    &resource.AtType,
		Access:    &resource.Access,
		Preserve:  &resource.Preserve,
		Publish:   &resource.Publish}
}

func loadParams(id string, params operations.DepositResourceParams) interface{} {
	// NOTE: This section will be replaced by DataUtils
	return map[string]interface{}{
		"id":        id,
		"attype":    params.Payload.AtType,
		"atcontext": params.Payload.AtContext,
		"access":    params.Payload.Access,
		"label":     params.Payload.Label,
		"preserve":  params.Payload.Preserve,
		"publish":   params.Payload.Publish,
		"sourceid":  params.Payload.SourceID,
	}
}
