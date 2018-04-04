package db

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/datautils"
)

// ReadVersion - return the resource with the particular id and version
func (database DynamodbDatabase) ReadVersion(id string, version *string) (*datautils.Resource, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
			"version": {
				S: version,
			},
		},
		TableName:      &database.Table,
		ConsistentRead: aws.Bool(true),
	}

	return database.query(params)
}

func (database DynamodbDatabase) query(params *dynamodb.GetItemInput) (*datautils.Resource, error) {
	resp, err := database.Connection.GetItem(params)
	if err != nil {
		return nil, err
	}
	return database.deserializeResponse(resp)
}

func (database DynamodbDatabase) deserializeResponse(resp *dynamodb.GetItemOutput) (*datautils.Resource, error) {
	if len(resp.Item) == 0 {
		return nil, errors.New("not found")
	}

	var json datautils.JSONObject
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &json); err != nil {
		return nil, err
	}

	return datautils.NewResource(json), nil
}
