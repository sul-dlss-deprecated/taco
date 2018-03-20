package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

// When we pass in an object that isn't a DRO, just let it pass through.
func TestSequenceValidatorNotApplicable(t *testing.T) {
	validator := NewSequenceValidator()
	err := validator.ValidateResource(testFileResource("whatever"))
	assert.Nil(t, err)
}

func TestSequenceIsASubset(t *testing.T) {
	validator := NewSequenceValidator()
	seq1 := map[string]interface{}{
		"members": []interface{}{"9", "10", "11"},
	}
	structural := map[string]interface{}{
		"hasMember":       []interface{}{"9", "10", "11"},
		"hasMemberOrders": []interface{}{seq1},
	}
	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-image.jsonld",
		"structural": structural}
	obj := datautils.NewResource(json)
	err := validator.ValidateResource(obj)
	assert.Nil(t, err)
}

func TestSequenceIsASuperset(t *testing.T) {
	validator := NewSequenceValidator()
	seq1 := map[string]interface{}{
		"members": []interface{}{"11", "9", "13"},
	}
	structural := map[string]interface{}{
		"hasMember":       []interface{}{"9", "10", "11"},
		"hasMemberOrders": []interface{}{seq1},
	}
	json := datautils.JSONObject{
		"@type":      "http://sdr.sul.stanford.edu/models/sdr3-image.jsonld",
		"structural": structural}
	obj := datautils.NewResource(json)
	err := validator.ValidateResource(obj)
	assert.NotNil(t, err)
}
