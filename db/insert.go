package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Insert(database *Database, params interface{}) error {
	row, err := dynamodbattribute.MarshalMap(params)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: aws.String(database.Table),
	}

	_, err = database.Connection.PutItem(input)

	if err != nil {
		return err
	}

	return nil

}
