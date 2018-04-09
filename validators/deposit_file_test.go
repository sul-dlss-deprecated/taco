package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestDepositFileIsValid(t *testing.T) {
	resource := &datautils.Resource{}
	err := NewDepositFileValidator(newMockRepository(nil)).ValidateResource(resource)
	assert.Nil(t, err)
}
