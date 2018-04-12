package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Insert create a row in dynamodb
func (database *DynamodbDatabase) Insert(resource *datautils.Resource) error {
	row, err := MarshalMap(resource.JSON)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: aws.String(database.Table),
	}

	_, err = database.Connection.PutItem(input)
	return err
}
