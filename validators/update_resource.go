package validators

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
	"github.com/sul-dlss-labs/taco/db"
)

// UpdateResourceValidator validates the update resource request
type UpdateResourceValidator struct {
	repository db.Database
	schema     *jsonschema.Schema
}

// NewUpdateResourceValidator creates a new instance of UpdateResourceValidator
func NewUpdateResourceValidator(repository db.Database, schemaPath string) *UpdateResourceValidator {
	schema, err := jsonschema.Compile(schemaPath)
	if err != nil {
		panic(err)
	}
	return &UpdateResourceValidator{repository: repository,
		schema: schema}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *UpdateResourceValidator) ValidateResource(body string) error {
	f := strings.NewReader(body)
	if err := d.schema.Validate(f); err != nil {
		return err
	}
	return nil
}
