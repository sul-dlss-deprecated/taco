package validators

import (
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type sequenceValidator struct{}

// NewSequenceValidator returns a validator that checks the structure metadata of any sequences
func NewSequenceValidator() ResourceValidator {
	return &sequenceValidator{}
}

// ValidateResource returns nil if the type is not a Collection or Object type.
// Otherwise, it checks and returns errors if any member of a sequence is not asserted in hasMember
func (d *sequenceValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {

	if !resource.IsCollection() && !resource.IsObject() {
		// Nothing to do here, this is not a Collection or Object
		return nil
	}

	// Load the member resource
	structuralMd := resource.Structural()
	if !structuralMd.HasKey("hasMemberOrders") {
		// Nothing to do here, there are no orders to validate
		return nil
	}
	if !structuralMd.HasKey("hasMember") {
		// Must have hasMember if they have hasMemberOrders
		return d.buildErrors("Resource has 'hasMemberOrders' but does not have 'hasMember'", "structural.hasMember")
	}
	memberOrders := structuralMd.GetA("hasMemberOrders")
	memberIds := structuralMd.GetA("hasMember").GetS()

	for _, order := range memberOrders.GetObj() {
		// check that order.members is a subset of memberIds
		membersOfOrder := order.GetA("members").GetS()
		if !subset(membersOfOrder, memberIds) {
			return d.buildErrors("'hasMemberOrders.members' must be a subset of 'hasMember'", "structural.hasMemberOrders.member")
		}
	}

	return nil
}

func (d *sequenceValidator) buildErrors(message string, pointer string) *models.ErrorResponseErrors {
	errors := models.ErrorResponseErrors{}
	source := &models.ErrorSource{Pointer: pointer}
	problem := &models.Error{
		Title:  "Validation Error",
		Detail: message,
		Source: source,
	}
	errors = append(errors, problem)
	return &errors
}

// subset returns true if the first array is completely
// contained in the second array. There must be at least
// the same number of duplicate values in second as there
// are in first.
func subset(first, second []string) bool {
	set := make(map[string]int)
	for _, value := range second {
		set[value]++
	}

	for _, value := range first {
		if count, found := set[value]; !found {
			return false
		} else if count < 1 {
			return false
		} else {
			set[value] = count - 1
		}
	}

	return true
}
