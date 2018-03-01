package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateResourceIsValid(t *testing.T) {
	err := NewUpdateResourceValidator(newMockRepository()).ValidateResource(testResource())
	assert.Nil(t, err)
}
