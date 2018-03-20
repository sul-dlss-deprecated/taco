package validators

import (
	"encoding/json"
	"strings"

	"github.com/santhosh-tekuri/jsonschema"
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// DepositResourceValidator validates the deposit resource request
type DepositResourceValidator struct {
	repository db.Database
	schema     *jsonschema.Schema
}

// NewDepositResourceValidator creates a new instance of DepositResourceValidator
func NewDepositResourceValidator(repository db.Database) ResourceValidator {
	files := []string{"DepositResource.json", "DepositCollection.json", "Sequence.json", "Agent.json", "DepositDRO.json", "DepositFileset.json", "DepositFile.json"}
	schema := BuildSchema("DepositResource.json", files)
	return &DepositResourceValidator{repository: repository, schema: schema}
}

// ValidateResource validates that a Resource models is semantically acceptable
func (d *DepositResourceValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	if err := d.schema.Validate(toReader(resource)); err != nil {
		return buildErrors(err.(*jsonschema.ValidationError))
	}
	return nil
}

func toReader(resource *datautils.Resource) *strings.Reader {
	json, err := json.Marshal(resource.JSON)
	if err != nil {
		panic(err)
	}
	return strings.NewReader(string(json[:]))
}
