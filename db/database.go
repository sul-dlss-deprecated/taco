package db

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Database struct {
	Connection *dynamodb.DynamoDB
	Table      string
}
