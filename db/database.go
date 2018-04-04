package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Database is a generic connection to a database.
type Database interface {
	Insert(*datautils.Resource) error
	Read(id string) (*datautils.Resource, error)
	DeleteByID(id string) error
	ReadVersion(id string, version *string) (*datautils.Resource, error)
}

// DynamodbDatabase Represents a connection to Dynamo
type DynamodbDatabase struct {
	Connection *dynamodb.DynamoDB
	Table      string
}

// Connect creates a dynamodb connection
func Connect(session *session.Session, dynamodbEndpoint string) *dynamodb.DynamoDB {
	dynamoConfig := &aws.Config{Endpoint: aws.String(dynamodbEndpoint)}
	return dynamodb.New(session, dynamoConfig)
}

func (database DynamodbDatabase) Read(id string) (*datautils.Resource, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:      &database.Table,
		ConsistentRead: aws.Bool(true),
	}
	return database.query(params)
}
