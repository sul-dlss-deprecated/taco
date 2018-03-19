package validators

import (
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
	"github.com/sul-dlss-labs/taco/db"
)

// DepositResourceValidator validates the deposit resource request
type DepositResourceValidator struct {
	repository db.Database
	schema     *jsonschema.Schema
}

// NewDepositResourceValidator creates a new instance of DepositResourceValidator
func NewDepositResourceValidator(repository db.Database) *DepositResourceValidator {
	files := []string{"DepositResource.json", "DepositCollection.json", "Sequence.json", "Agent.json", "DepositDRO.json", "DepositFileset.json", "DepositFile.json"}
	schema := BuildSchema("DepositResource.json", files)
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
