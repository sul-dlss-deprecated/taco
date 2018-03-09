package resource

import (
	"log"
	"reflect"

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
		log.Printf("In ERRROR")
		log.Printf("%#v", row)
		return err
	} else {
		log.Printf("SUCCESSFUL STRUCT")
	}
	log.Printf("%#v", row)
	log.Printf("Saved %s to dynamodb", id)

	return nil
}

func loadParams(id string, params operations.DepositResourceParams) interface{} {
	paramsValue := reflect.ValueOf(params.Payload)
	paramsValueFields := reflect.ValueOf(params.Payload).Elem()
	paramsIndirect := reflect.Indirect(paramsValue)
	paramsType := paramsIndirect.Type()
	log.Printf("ParamsType: %v", paramsType)
	//	log.Printf("ParamsType: %v", paramsType)
	for i := 0; i < paramsType.NumField(); i++ {
		paramField := paramsType.Field(i)
		paramValue := paramsValueFields.Field(i)

		if reflect.Ptr != paramValue.Kind() {
			if paramValue.Len() > 0 {
				// this is apparently a slice
			}
		} else {
			log.Printf("FieldName: %v - FieldValue: %v", paramField.Name, paramValue.Elem().Interface())
		}
		//		paramsField := paramsType.Field(i)
		//		log.Printf("Parmas Field: %s", paramsField.Name)
	}

	//	var m map[string]interface{}
	//	err := json.Unmarshal(params.Payload, &m)

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
