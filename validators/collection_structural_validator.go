package validators

import (
	"fmt"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type collectionStructuralValidator struct {
	repository db.Database
}

// NewCollectionStructuralValidator returns a validator that checks the structure metadata of an Object (DRO)
func NewCollectionStructuralValidator(repository db.Database) ResourceValidator {
	return &collectionStructuralValidator{repository: repository}
}

// ValidateResource returns nil if the type is not a Collection type.
// Otherwise, it checks and returns errors if each member is not an object
func (d *collectionStructuralValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {

	if !resource.IsCollection() {
		// Nothing to do here, this is not a Collection
		return nil
	}

	// Load the member resource
	structuralMd := resource.Structural()
	if !structuralMd.HasKey("hasMember") {
		// Nothing to do here, there are no members to validate
		return nil
	}

	memberIds := structuralMd.GetA("hasMember").GetS()
	for _, memberID := range memberIds {
		member, err := d.repository.RetrieveLatest(memberID)
		if err != nil {
			return d.buildErrors(fmt.Sprintf("Unable to find member %s", memberID))
		}
		if !member.IsObject() {
			return d.buildErrors(fmt.Sprintf("Collection members must be object types. %s is a %s", member.ID(), member.Type()))
		}
	}

	return nil
}

func (d *collectionStructuralValidator) buildErrors(message string) *models.ErrorResponseErrors {
	errors := models.ErrorResponseErrors{}
	source := &models.ErrorSource{Pointer: "structural.hasMember"}
	problem := &models.Error{
		Title:  "Validation Error",
		Detail: message,
		Source: source,
	}
	errors = append(errors, problem)
	return &errors
}
