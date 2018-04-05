package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// UpdateString - Updates a string value in the dynamo repository for this resource
func (h DynamodbDatabase) UpdateString(resourceID string, field string, value string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(value),
			},
		},
		TableName: &h.Table,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(resourceID),
			},
		},
		UpdateExpression: aws.String("set " + field + " = :r"),
	}

	_, err := h.Connection.UpdateItem(input)

	if err != nil {
		return err
	}

	return nil
}

// UpdateBool - Updates a boolean value in the dynamo repository for this resource
func (h DynamodbDatabase) UpdateBool(resourceID string, field string, value bool) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				BOOL: aws.Bool(value),
			},
		},
		TableName: &h.Table,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(resourceID),
			},
		},
		UpdateExpression: aws.String("set " + field + " = :r"),
	}

	_, err := h.Connection.UpdateItem(input)

	if err != nil {
		return err
	}

	return nil
}
