package db

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Database is a generic connection to a database.
type Database interface {
	Insert(*datautils.Resource) error
	DeleteAllVersions(externalID string) error
	RetrieveVersion(externalID string, version *string) (*datautils.Resource, error)
	RetrieveLatest(externalID string) (*datautils.Resource, error)
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

func (d *DynamodbDatabase) query(params *dynamodb.QueryInput) (*datautils.Resource, error) {
	resp, err := d.Connection.Query(params)

	if err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, errors.New("not found")
	}

	return respToResource(resp.Items[0])
}

func respToResource(item map[string]*dynamodb.AttributeValue) (*datautils.Resource, error) {
	var json datautils.JSONObject
	if err := UnmarshalMap(item, &json); err != nil {
		return nil, err
	}

	// UnmarshalMap coerces all numbers in AWS to float64.
	// Force version to be an integer
	json["version"] = int(json["version"].(float64))

	return datautils.NewResource(json), nil
}
