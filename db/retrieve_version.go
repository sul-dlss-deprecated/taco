package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/datautils"
)

// RetrieveVersion - return the resource with the particular externalID and version
func (d *DynamodbDatabase) RetrieveVersion(externalID string, version *string) (*datautils.Resource, error) {

	params := &dynamodb.QueryInput{
		KeyConditions: map[string]*dynamodb.Condition{
			"externalIdentifier": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(externalID),
					},
				},
			},
			"version": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						N: version,
					},
				},
			},
		},
		IndexName:        aws.String("ResourceByExternalIDAndVersion"),
		Limit:            aws.Int64(1),
		ScanIndexForward: aws.Bool(true),
		TableName:        &d.Table,
	}

	return d.query(params)
}
