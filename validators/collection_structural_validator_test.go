package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func testCollectionResource(members []string) *datautils.Resource {
	structural := map[string]interface{}{}
	if len(members) != 0 {
		memberIDs := make([]interface{}, len(members))
		for i, v := range members {
			memberIDs[i] = v
		}
		structural["hasMember"] = memberIDs
	}

	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-curated-collection.jsonld",
		"structural": structural,
	}
	return datautils.NewResource(json)
}

// When we pass in an object that isn't a Collection, just let it pass through.
func TestCollectionStructuralValidatorNotCollection(t *testing.T) {
	validator := NewCollectionStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(testFilesetResource(""))
	assert.Nil(t, err)
}

func TestCollectionStructuralValidatorValid(t *testing.T) {
	col := testCollectionResource([]string{"object2"})
	obj := testObjectResource([]string{})
	validator := NewCollectionStructuralValidator(newMockRepository(obj))
	err := validator.ValidateResource(col)
	assert.Nil(t, err)
}

// When the member is not an Object
func TestCollectionStructuralValidatorWrongType(t *testing.T) {
	collection := testResource("bs646cd8717.json")
	validator := NewCollectionStructuralValidator(newMockRepository(collection))
	col := testCollectionResource([]string{collection.ID()})
	err := validator.ValidateResource(col)
	assert.NotNil(t, err)
}

// when we can't validate the reference to the member
func TestCollectionStructuralValidatorMemberNotFound(t *testing.T) {
	validator := NewCollectionStructuralValidator(newMockRepository(nil))
	col := testCollectionResource([]string{"NotfindableID"})
	err := validator.ValidateResource(col)
	assert.NotNil(t, err)
}

// when there is no member assertion in the metadata
func TestCollectionStructuralValidatorNoMemberAssertion(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-collection.jsonld",
		"structural": map[string]interface{}{}}
	obj := datautils.NewResource(json)
	validator := NewCollectionStructuralValidator(newMockRepository(nil))
	err := validator.ValidateResource(obj)
	assert.Nil(t, err)
}
