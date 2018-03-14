package db

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// Database is a generic connection to a database.
type Database interface {
	Insert(interface{}) error
	Update(interface{}) error
	Read(id string) (*models.Resource, error)
}

// DynamodbDatabase Represents a connection to Dynamo
type DynamodbDatabase struct {
	Connection *dynamodb.DynamoDB
	Table      string
}

func Connect(session *session.Session, dynamodbEndpoint string) *dynamodb.DynamoDB {
	dynamoConfig := &aws.Config{Endpoint: aws.String(dynamodbEndpoint)}
	return dynamodb.New(session, dynamoConfig)
}

func (database DynamodbDatabase) Read(id string) (*models.Resource, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:      &database.Table,
		ConsistentRead: aws.Bool(true),
	}
	resp, err := database.Connection.GetItem(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var resource *models.Resource
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &resource); err != nil {
		log.Println(err)
		return nil, err
	}

	if resource.ID == "" {
		return nil, errors.New("not found")
	}
	return resource, nil
}
