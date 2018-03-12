package validators

import (
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// DepositResourceValidator validates the deposit resource request
type DepositResourceValidator struct {
	repository db.Database
}

// NewDepositResourceValidator creates a new instance of DepositResourceValidator
func NewDepositResourceValidator(repository db.Database) *DepositResourceValidator {
	return &DepositResourceValidator{repository: repository}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *DepositResourceValidator) ValidateResource(resource *models.Resource) error {
	// TODO: Add checks here
	return nil
}
