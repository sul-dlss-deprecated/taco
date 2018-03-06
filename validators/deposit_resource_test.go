package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepositResourceIsValid(t *testing.T) {
	err := NewDepositResourceValidator(newMockRepository()).ValidateResource(testDepositResource())
	assert.Nil(t, err)
}
