package validators

import (
	"mime/multipart"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDepositFileIsValid(t *testing.T) {
	err := NewDepositFileValidator(newMockRepository()).ValidateResource(&multipart.FileHeader{})
	assert.Nil(t, err)
}
