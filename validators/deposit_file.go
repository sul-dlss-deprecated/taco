package validators

import (
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// DepositFileValidator validates the deposit file request
type DepositFileValidator struct {
	repository db.Database
}

// NewDepositFileValidator creates a new instance of DepositFileValidator
func NewDepositFileValidator(repository db.Database) ResourceValidator {
	return &DepositFileValidator{repository: repository}
}

// ValidateResource validates that the file resource is acceptable
func (d *DepositFileValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	// TODO: Add checks here
	return nil
}
