package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
)

var db *dynamodb.DynamoDB

// NewConnection creates a new connection to Dynamo using our config
func NewConnection(c *viper.Viper) *dynamodb.DynamoDB {
	return dynamodb.New(session.New(&aws.Config{
		Region:      aws.String(c.GetString("db.region")),
		Credentials: credentials.NewEnvCredentials(),
		Endpoint:    aws.String(c.GetString("db.endpoint")),
		DisableSSL:  aws.Bool(c.GetBool("db.disable_ssl")),
	}))
}
