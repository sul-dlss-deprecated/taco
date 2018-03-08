package resource

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

func Create(database *dynamodb.DynamoDB, id string, params operations.DepositResourceParams) error {
	row, err := dynamodbattribute.MarshalMap(loadParams(id, params))

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: aws.String("resources"),
	}

	_, err = database.PutItem(input)

	if err != nil {
		return err
	}
	log.Printf("Saved %s to dynamodb", id)

	return nil
}

func loadParams(id string, params operations.DepositResourceParams) *Resource {
	return &Resource{
		ID:        id,
		AtType:    *params.Payload.AtType,
		AtContext: *params.Payload.AtContext,
		Access:    *params.Payload.Access,
		Label:     *params.Payload.Label,
		Preserve:  *params.Payload.Preserve,
		Publish:   *params.Payload.Publish,
		SourceID:  params.Payload.SourceID,
	}
}
