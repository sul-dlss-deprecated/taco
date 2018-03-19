package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepositResourceIsValid(t *testing.T) {
	depositValidator := NewDepositResourceValidator(newMockRepository())
	err := depositValidator.ValidateResource(testDepositResource())
	assert.Nil(t, err)
}
