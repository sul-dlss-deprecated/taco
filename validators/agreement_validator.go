package validators

import (
	"fmt"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type agreementValidator struct {
	repository db.Database
}

// NewAgreementValidator returns a validator that checks that the agreement exists and is valid
func NewAgreementValidator(repository db.Database) ResourceValidator {
	return &agreementValidator{repository: repository}
}

// ValidateResource returns nil if the type is not a Collection or DRO type.
// Otherwise, it checks and returns errors if the hasAgreement assertion is not an Agreement
func (d *agreementValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {

	if !resource.IsCollection() && !resource.IsObject() {
		// Nothing to do here, this is not a Collection or DRO
		return nil
	}

	// Load the member resource
	structuralMd := resource.Structural()
	if !structuralMd.HasKey("hasAgreement") {
		// Nothing to do here, there is no agreement to validate
		return nil
	}

	agreementID := structuralMd.GetS("hasAgreement")
	agreement, err := d.repository.RetrieveLatest(agreementID)
	if err != nil {
		return d.buildErrors(fmt.Sprintf("Unable to find agreement %s", agreementID))
	}

	if !agreement.IsAgreement() {
		return d.buildErrors(fmt.Sprintf("Agreements must be Agreement type. %s is a %s", agreementID, agreement.Type()))
	}

	return nil
}

func (d *agreementValidator) buildErrors(message string) *models.ErrorResponseErrors {
	errors := models.ErrorResponseErrors{}
	source := &models.ErrorSource{Pointer: "structural.hasAgreement"}
	problem := &models.Error{
		Title:  "Validation Error",
		Detail: message,
		Source: source,
	}
	errors = append(errors, problem)
	return &errors
}
