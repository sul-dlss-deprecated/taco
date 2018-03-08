package validators

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
)

// DepositResourceValidator validates the deposit resource request
type DepositResourceValidator struct {
	repository db.Repository
	schema     *jsonschema.Schema
}

// NewDepositResourceValidator creates a new instance of DepositResourceValidator
func NewDepositResourceValidator(repository db.Repository, schemaPath string) *DepositResourceValidator {
	schema, err := jsonschema.Compile(schemaPath)
	if err != nil {
		panic(err)
	}
	return &DepositResourceValidator{repository: repository, schema: schema}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *DepositResourceValidator) ValidateResource(body string) error {
	f := strings.NewReader(body)
	if err := d.schema.Validate(f); err != nil {
		return err
	}
	return nil
}
