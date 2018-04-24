package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DeleteByID -- given an identifier, remove the resource
func (h DynamodbDatabase) DeleteByID(tacoIdentifier string) error {
	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"tacoIdentifier": {
				S: aws.String(tacoIdentifier),
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
