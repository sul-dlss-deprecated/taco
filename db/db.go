package db

import (
	"github.com/sul-dlss-labs/taco/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var db *dynamodb.DynamoDB

// NewConnection creates a new connection to Dynamo using our config
func NewConnection() *dynamodb.DynamoDB {
	config := config.NewConfig()
	return dynamodb.New(session.New(&aws.Config{
		Region:      aws.String(config.AWS_Region),
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(config.DB_Endpoint),
		DisableSSL:  aws.Bool(config.Disable_SSL),
	}))
}
