package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/config"
)

var db *dynamodb.DynamoDB

// NewConnection creates a new connection to Dynamo using our config
func NewConnection(config *config.Config, sess *session.Session) *dynamodb.DynamoDB {
	dynamoConfig := &aws.Config{Endpoint: aws.String(config.DynamodbEndpoint)}
	return dynamodb.New(sess, dynamoConfig)
}
