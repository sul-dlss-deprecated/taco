package validators

import (
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/persistence"
)

// DepositResourceValidator validates the deposit resource request
type DepositResourceValidator struct {
	repository persistence.Repository
}

// NewDepositResourceValidator creates a new instance of DepositResourceValidator
func NewDepositResourceValidator(repository persistence.Repository) *DepositResourceValidator {
	return &DepositResourceValidator{repository: repository}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *DepositResourceValidator) ValidateResource(resource *models.DepositResource) error {
	// TODO: Add checks here
	return nil
}
