package validators

import (
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/sul-dlss-labs/taco/generated/models"
)

func buildErrors(validationError *jsonschema.ValidationError) *models.ErrorResponseErrors {
	// The top level error will be: I[#] S[#] doesn't validate with "DepositResource.json#"
	// we really want the causes of why it didn't validate.
	causes := validationError.Causes
	errors := models.ErrorResponseErrors{}
	for _, cause := range causes {
		source := &models.ErrorSource{Pointer: cause.InstancePtr}
		problem := &models.Error{Title: "Validation Error", Detail: cause.Error(), Source: source}
		errors = append(errors, problem)
	}
	return &errors
}
