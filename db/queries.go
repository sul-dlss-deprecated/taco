package db

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// AnyWithDedupeIdentifier -- are there any resources with the given dedupeIdentifier?
func (h DynamodbDatabase) AnyWithDedupeIdentifier(dedupeIdentifier string) (bool, error) {
	params := &dynamodb.QueryInput{
		TableName: aws.String(h.Table),
		IndexName: aws.String(h.IndexName),
		KeyConditions: map[string]*dynamodb.Condition{
			"dedupeIdentifier": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(dedupeIdentifier),
					},
				},
			},
		},
	}
	resp, err := h.Connection.Query(params)

	if err != nil {
		log.Println(err)
		return false, err
	}
	return (*resp.Count > 0), nil
}
