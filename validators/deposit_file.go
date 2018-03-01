package validators

import (
	"mime/multipart"

	"github.com/sul-dlss-labs/taco/persistence"
)

// DepositFileValidator validates the deposit file request
type DepositFileValidator struct {
	repository persistence.Repository
}

// NewDepositFileValidator creates a new instance of DepositFileValidator
func NewDepositFileValidator(repository persistence.Repository) *DepositFileValidator {
	return &DepositFileValidator{repository: repository}
}

// ValidateResource validates that a headers are acceptable
func (d *DepositFileValidator) ValidateResource(resource *multipart.FileHeader) error {
	// TODO: Add checks here
	return nil
}
