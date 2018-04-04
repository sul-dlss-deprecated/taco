package validators

import "github.com/sul-dlss-labs/taco/generated/models"

// ResourceValidator validates that a request body is acceptable
type ResourceValidator interface {
	ValidateResource(body string) *models.ErrorResponseErrors
}
