package validators

import (
	"fmt"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type fileStructuralValidator struct {
	repository db.Database
}

// NewFileStructuralValidator returns a validator that checks the structure metadata of a File
func NewFileStructuralValidator(repository db.Database) ResourceValidator {
	return &fileStructuralValidator{repository: repository}
}

// ValidateResource returns nil if the type is not a file type.
// Otherwise, it checks and returns errors if isContainedBy is not a fileset
func (d *fileStructuralValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	if !resource.IsFile() {
		// Nothing to do here, this is not a File
		return nil
	}
	// Load the containing resource
	structuralMd := resource.Structural()
	if !structuralMd.HasKey("isContainedBy") {
		return d.buildErrors("A file must be contained by a Fileset")
	}

	containedBy := structuralMd.GetS("isContainedBy")
	containing, err := d.repository.RetrieveLatest(containedBy)
	if err != nil {
		return d.buildErrors(fmt.Sprintf("Unable to find container %s", containedBy))
	}
	if containing.IsFileset() {
		return nil
	}

	return d.buildErrors(fmt.Sprintf("Files may only be contained by Filesets.  You specified %s, which is a %s", containedBy, containing.Type()))
}

func (d *fileStructuralValidator) buildErrors(message string) *models.ErrorResponseErrors {
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
