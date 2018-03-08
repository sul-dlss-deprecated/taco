package persistence

import (
	"github.com/sul-dlss-labs/taco/generated/models"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// ToPersistable converts the model Resource into a persisable object
func NewResourceWithId(resourceID string, request *models.Resource) *Resource {
	resource := NewResource()
	resource.MarshalMap(request)
	resource.PutS("id", resourceID)
	return resource
}

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	data map[string]*dynamodb.AttributeValue
}

// NewResource returns a new instance of the resource
func NewResource() *Resource {
	return &Resource{data: map[string]*dynamodb.AttributeValue{}}
}

// PutS add a string value to the Resource
func (d *Resource) PutS(key string, value string) {
	d.data[key] = &dynamodb.AttributeValue{S: &value}
}

// GetS return a string value from the Resource
func (d *Resource) GetS(key string) string {
	return *d.data[key].S
}

// MarshalMap deserializes the provided struct into this resource
func (d *Resource) MarshalMap(request interface{}) {
	row, err := dynamodbattribute.MarshalMap(request)
	if err != nil {
		panic(err)
	}
	d.data = row
}

// Map return the data as a map of string -> dynamodb.AttributeValue
func (d *Resource) Map() map[string]*dynamodb.AttributeValue {
	return d.data
}

// ToPersistable converts the model Resource into a persisable object
func ToPersistable(resourceID string, request *models.Resource) *Resource {
	resource := NewResource()
	resource.MarshalMap(request)
	resource.PutS(PrimaryKey, resourceID)
	return resource
}
