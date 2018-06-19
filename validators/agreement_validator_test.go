package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func testResourceWithAgreement(agreementID string) *datautils.Resource {
	json := datautils.JSONObject{
		"@type": datautils.AgreementType,
		"structural": map[string]interface{}{
			"hasAgreement": agreementID,
		}}
	return datautils.NewResource(json)
}

// When we pass in an object that isn't an object or collection, just let it pass through.
func TestAgreementValidatorNotDRO(t *testing.T) {
	validator := NewAgreementValidator(newMockRepository(nil))
	err := validator.ValidateResource(testFilesetResource("999999"))
	assert.Nil(t, err)
}

func TestAgreementValidatorValid(t *testing.T) {
	json := datautils.JSONObject{
		"@type": datautils.AgreementType,
	}
	agreement := datautils.NewResource(json)
	object := testResourceWithAgreement("object1")
	validator := NewAgreementValidator(newMockRepository(agreement))
	err := validator.ValidateResource(object)
	assert.Nil(t, err)
}

// When the agreement object is not an agreement type
func TestAgreementValidatorWrongType(t *testing.T) {
	notAgreement := testResource("bs646cd8717.json")
	json := datautils.JSONObject{
		"@type":      datautils.ObjectTypes[0],
		"structural": map[string]interface{}{"hasAgreement": "99999"}}
	obj := datautils.NewResource(json)
	validator := NewAgreementValidator(newMockRepository(notAgreement))
	err := validator.ValidateResource(obj)
	assert.NotNil(t, err)
}

// when we can't validate the reference to the container
func TestAgreementValidatorAgreementNotFound(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      datautils.ObjectTypes[0],
		"structural": map[string]interface{}{"hasAgreement": "99999"}}
	obj := datautils.NewResource(json)
	validator := NewAgreementValidator(newMockRepository(nil))
	err := validator.ValidateResource(obj)
	assert.NotNil(t, err)
}

// when there is no hasAgreement assertion in the metadata
func TestAgreementValidatorNoAgreementAssertion(t *testing.T) {
	json := datautils.JSONObject{
		"@type":      datautils.ObjectTypes[0],
		"structural": map[string]interface{}{}}
	obj := datautils.NewResource(json)
	validator := NewAgreementValidator(newMockRepository(nil))
	err := validator.ValidateResource(obj)
	assert.Nil(t, err)
}
