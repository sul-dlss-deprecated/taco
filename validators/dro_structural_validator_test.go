package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func testObjectResource(members []string) *datautils.Resource {
	structural := map[string]interface{}{}
	if len(members) != 0 {
		memberIDs := make([]interface{}, len(members))
		for i, v := range members {
			memberIDs[i] = v
		}
		structural["hasMember"] = memberIDs
	}

	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-photograph.jsonld",
		"structural": structural,
	}
	return datautils.NewResource(json)
}

// When we pass in an object that isn't a DRO, just let it pass through.
func TestDROStructuralValidatorNotObject(t *testing.T) {
	validator := NewDROStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testResource("bs646cd8717.json"))
	assert.Nil(t, err)
}

func TestDROStructuralValidatorValid(t *testing.T) {
	object1 := testObjectResource([]string{"object2"})
	object2 := testObjectResource([]string{})
	validator := NewDROStructuralValidator(newMockRepository(object2))
	err := validator.ValidateResource(object1)
	assert.Nil(t, err)
}

// When the member is not an Object
func TestDROStructuralValidatorWrongType(t *testing.T) {
	collection := testResource("bs646cd8717.json")
	validator := NewDROStructuralValidator(newMockRepository(collection))
	obj := testObjectResource([]string{collection.ID()})
	err := validator.ValidateResource(obj)
	assert.NotNil(t, err)
}

// when we can't validate the reference to the member
func TestDROStructuralValidatorMemberNotFound(t *testing.T) {
	validator := NewDROStructuralValidator(newMockRepository(nil))
	obj := testObjectResource([]string{"NotfindableID"})
	err := validator.ValidateResource(obj)
	assert.NotNil(t, err)
}

// when there is no member assertion in the metadata
func TestDROStructuralValidatorNoMemberAssertion(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-image.jsonld",
		"structural": map[string]interface{}{}}
	obj := datautils.NewResource(json)
	validator := NewDROStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(obj)
	assert.Nil(t, err)
}
