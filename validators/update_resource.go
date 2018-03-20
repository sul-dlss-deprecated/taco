package validators

import (
	"github.com/santhosh-tekuri/jsonschema"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// UpdateResourceValidator validates the update resource request
type UpdateResourceValidator struct {
	repository db.Database
	schema     *jsonschema.Schema
}

// NewUpdateResourceValidator creates a new instance of UpdateResourceValidator
func NewUpdateResourceValidator(repository db.Database) ResourceValidator {
	files := []string{"Resource.json", "Collection.json", "Sequence.json", "Agent.json", "DRO.json", "Fileset.json", "File.json"}
	schema := BuildSchema("Resource.json", files)
	return &UpdateResourceValidator{repository: repository,
		schema: schema}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *UpdateResourceValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	if err := d.schema.Validate(toReader(resource)); err != nil {
		return buildErrors(err.(*jsonschema.ValidationError))
	}
	return nil
}
