package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	updateValidator = NewUpdateResourceValidator(newMockRepository(nil))
)

func TestUpdateResourceIsValid(t *testing.T) {
	resource := testResource("bs646cd8717.json")
	err := updateValidator.ValidateResource(resource)
	assert.Nil(t, err)
}

func TestUpdateResourceIsInvalid(t *testing.T) {
	resource := testResource("invalid-bq582kh2487.json")
	err := updateValidator.ValidateResource(resource)
	assert.Equal(t, "Validation Error", (*err)[0].Title)
}
