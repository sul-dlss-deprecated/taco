package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// deleteByID -- given an identifier, remove the resource
func (h DynamodbDatabase) deleteByID(id string) error {
	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"tacoIdentifier": {
				S: aws.String(id),
			},
		},
		TableName: aws.String(h.Table),
	}
	_, err := h.Connection.DeleteItem(params)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAllVersions removes all versions with the given external id
func (h *DynamodbDatabase) DeleteAllVersions(externalID string) error {
	resource, err := h.RetrieveLatest(externalID)
	if err != nil {
		panic(err)
	}
	// Delete all versions of the resource
	for resource != nil {
		if err = h.deleteByID(resource.ID()); err != nil {
			panic(err)
		}
		resource, err = h.RetrieveLatest(externalID)
		if _, ok := err.(*RecordNotFound); !ok {
			panic(err)
		}
	}
	return nil
}
