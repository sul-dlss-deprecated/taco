package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	depositValidator = NewDepositResourceValidator(newMockRepository(nil))
)

func TestDepositResourceIsValid(t *testing.T) {
	resource := testResource("create-bs646cd8717.json")
	err := depositValidator.ValidateResource(resource)
	assert.Nil(t, err)
}

func TestDepositResourceIsInvalid(t *testing.T) {
	resource := testResource("invalid-bq582kh2487.json")
	err := depositValidator.ValidateResource(resource)
	assert.Equal(t, "Validation Error", (*err)[0].Title)
}
