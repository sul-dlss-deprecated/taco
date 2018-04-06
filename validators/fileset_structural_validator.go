package validators

import (
	"fmt"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type filesetStructuralValidator struct {
	repository db.Database
}

// NewFilesetStructuralValidator returns a validator that checks the structure metadata of a Fileset
func NewFilesetStructuralValidator(repository db.Database) ResourceValidator {
	return &filesetStructuralValidator{repository: repository}
}

// ValidateResource returns nil if the type is not a fileset type.
// Otherwise, it checks and returns errors if isContainedBy is not an object
func (d *filesetStructuralValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {

	if !resource.IsFileset() {
		// Nothing to do here, this is not a Fileset
		return nil
	}
	// Load the containing resource
	structuralMd := resource.Structural()
	if !structuralMd.HasKey("isContainedBy") {
		return d.buildErrors("A Fileset must be contained by an Object")
	}

	containedBy := structuralMd.GetS("isContainedBy")
	containing, err := d.repository.RetrieveLatest(containedBy)
	if err != nil {
		return d.buildErrors(fmt.Sprintf("Unable to find container %s", containedBy))
	}
	if containing.IsObject() {
		return nil
	}
	return d.buildErrors(fmt.Sprintf("Filesets may only be contained by Objects.  You specified %s, which is a %s", containedBy, containing.Type()))
}

func (d *filesetStructuralValidator) buildErrors(message string) *models.ErrorResponseErrors {
	errors := models.ErrorResponseErrors{}
	source := &models.ErrorSource{Pointer: "structural.isContainedBy"}
	problem := &models.Error{
		Title:  "Validation Error",
		Detail: message,
		Source: source,
	}
	errors = append(errors, problem)
	return &errors
}
