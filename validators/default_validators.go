package validators

import (
	"github.com/sul-dlss-labs/taco/db"
)

// NewDepositValidator builds the validator for deposit resource requests
func NewDepositValidator(database db.Database) ResourceValidator {
	return NewCompositeResourceValidator(
		[]ResourceValidator{
			NewDepositResourceValidator(database),
			structuralValidator(database),
		})
}

// NewUpdateValidator builds the validator for update resource requests
func NewUpdateValidator(database db.Database) ResourceValidator {
	return NewCompositeResourceValidator(
		[]ResourceValidator{
			NewUpdateResourceValidator(database),
			structuralValidator(database),
		})
}

// NewFileValidator builds the validator for deposit file requests
func NewFileValidator(database db.Database) ResourceValidator {
	return NewCompositeResourceValidator(
		[]ResourceValidator{
			NewDepositFileValidator(database),
			structuralValidator(database),
		})
}

// Builds the validator for structural validations.
// This is suitable for both create and update requests
func structuralValidator(database db.Database) ResourceValidator {
	return NewCompositeResourceValidator(
		[]ResourceValidator{
			NewFileStructuralValidator(database),
			NewFilesetStructuralValidator(database),
			NewDROStructuralValidator(database),
			NewCollectionStructuralValidator(database),
			NewAgreementValidator(database),
			NewSequenceValidator(),
		})
}
