package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/sul-dlss-labs/taco/datautils"
)

// RetrieveLatest - return the most recent version of the resource
// based on the max version for the identifier.
// Caveat: DynamoDB limits query results to 1MB, so really big metadata will
// break this
func (d *DynamodbDatabase) RetrieveLatest(externalID string) (*datautils.Resource, error) {
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
		},
		IndexName:        aws.String("ResourceByExternalIDAndVersion"),
		Limit:            aws.Int64(1),
		ScanIndexForward: aws.Bool(false), // Sort most recent version first.
		TableName:        &d.Table,
	}

	return d.query(params)
}
