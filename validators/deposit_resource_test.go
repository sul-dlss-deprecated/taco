package validators

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepositResourceIsValid(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps/DepositResource.json")
	err := NewDepositResourceValidator(newMockRepository(), schemaPath).ValidateResource(testDepositResource())
	assert.Nil(t, err)
}
