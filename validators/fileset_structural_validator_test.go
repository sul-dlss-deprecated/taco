package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func testFilesetResource(containerID string) *datautils.Resource {
	json := datautils.JSONObject{
		"@type": datautils.FilesetType,
		"structural": map[string]interface{}{
			"isContainedBy": containerID,
		}}
	return datautils.NewResource(json)
}

// When we pass in an object that isn't a fileset, just let it pass through.
func TestFilesetStructuralValidatorNotFileset(t *testing.T) {
	validator := NewFilesetStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testResource("bs646cd8717.json"))
	assert.Nil(t, err)
}

func TestFilesetStructuralValidatorValid(t *testing.T) {
	object := testObjectResource([]string{})
	fileset := testFilesetResource("object1")
	validator := NewFilesetStructuralValidator(newMockRepository(object))
	err := validator.ValidateResource(fileset)
	assert.Nil(t, err)
}

// When the containing object is not an Object
func TestFilesetStructuralValidatorWrongType(t *testing.T) {
	collection := testResource("bs646cd8717.json")
	validator := NewFilesetStructuralValidator(newMockRepository(collection))
	err := validator.ValidateResource(testFilesetResource(collection.ID()))
	assert.NotNil(t, err)
}

// when we can't validate the reference to the container
func TestFilesetStructuralValidatorContainerNotFound(t *testing.T) {
	validator := NewFilesetStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testFilesetResource("NotfindableID"))
	assert.NotNil(t, err)
}

// when there is no container assertion in the metadata
func TestFilesetStructuralValidatorNoContainedAssertion(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      datautils.FilesetType,
		"structural": map[string]interface{}{}}
	fileset := datautils.NewResource(json)
	validator := NewFilesetStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(fileset)
	assert.NotNil(t, err)
}
