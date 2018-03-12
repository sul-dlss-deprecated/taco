package validators

import (
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// UpdateResourceValidator validates the update resource request
type UpdateResourceValidator struct {
	repository db.Database
}

// NewUpdateResourceValidator creates a new instance of UpdateResourceValidator
func NewUpdateResourceValidator(repository db.Database) *UpdateResourceValidator {
	return &UpdateResourceValidator{repository: repository}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *UpdateResourceValidator) ValidateResource(resource *models.Resource) error {
	// TODO: Add checks here
	return nil
}
