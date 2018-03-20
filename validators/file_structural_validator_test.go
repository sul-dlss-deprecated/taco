package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func testFileResource(containerID string) *datautils.Resource {
	json := datautils.JSONObject{
		"@type": "http://sdr.sul.stanford.edu/models/sdr3-file.jsonld",
		"structural": map[string]interface{}{
			"isContainedBy": containerID,
		}}
	return datautils.NewResource(json)
}

// When we pass in an object that isn't a file, just let it pass through.
func TestFileStructuralValidatorNotFile(t *testing.T) {
	validator := NewFileStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testResource("bs646cd8717.json"))
	assert.Nil(t, err)
}

func TestFileStructuralValidatorValid(t *testing.T) {
	fileset := testFilesetResource("")
	file := testFileResource("fileset1")
	validator := NewFilesetStructuralValidator(newMockRepository(fileset))
	err := validator.ValidateResource(file)
	assert.Nil(t, err)
}

func TestFileStructuralValidatorWrongType(t *testing.T) {
	collection := testResource("bs646cd8717.json")
	validator := NewFileStructuralValidator(newMockRepository(collection))
	err := validator.ValidateResource(testFileResource(collection.ID()))
	assert.NotNil(t, err)
}

func TestFileStructuralValidatorContainerNotFound(t *testing.T) {
	validator := NewFileStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testFileResource("NotfindableID"))
	assert.NotNil(t, err)
}

func TestFileStructuralValidatorNoContainedAssertion(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-file.jsonld",
		"structural": map[string]interface{}{}}
	file := datautils.NewResource(json)
	validator := NewFileStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(file)
	assert.NotNil(t, err)
}
