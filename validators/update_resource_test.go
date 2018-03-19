package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateResourceIsValid(t *testing.T) {
	updateValidator := NewUpdateResourceValidator(newMockRepository())
	err := updateValidator.ValidateResource(testResource())
	assert.Nil(t, err)
}
