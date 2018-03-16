package validators

import (
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateResourceIsValid(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	schemaPath := path.Join(path.Dir(filename), "../maps")
	err := NewUpdateResourceValidator(newMockRepository(), schemaPath).ValidateResource(testResource())
	assert.Nil(t, err)
}
