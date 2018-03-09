package resource

import (
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

	return nil
}

func loadParams(id string, params operations.DepositResourceParams) interface{} {
	// NOTE: This section will be replaced by DataUtils
	return map[string]interface{}{
		"id":        id,
		"attype":    params.Payload.AtType,
		"atcontext": params.Payload.AtContext,
		"access":    params.Payload.Access,
		"label":     params.Payload.Label,
		"preserve":  params.Payload.Preserve,
		"publish":   params.Payload.Publish,
		"sourceid":  params.Payload.SourceID,
	}
}
