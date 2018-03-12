package validators

import (
	"mime/multipart"

	"github.com/sul-dlss-labs/taco/db"
)

// DepositFileValidator validates the deposit file request
type DepositFileValidator struct {
	repository db.Database
}

// NewDepositFileValidator creates a new instance of DepositFileValidator
func NewDepositFileValidator(repository db.Database) *DepositFileValidator {
	return &DepositFileValidator{repository: repository}
}

// ValidateResource validates that a headers are acceptable
func (d *DepositFileValidator) ValidateResource(resource *multipart.FileHeader) error {
	// TODO: Add checks here
	return nil
}
