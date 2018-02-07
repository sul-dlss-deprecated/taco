package db

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DeleteByID -- given an identifier, remove the resource
func (h DynamodbDatabase) DeleteByID(tacoIdentifier string) error {
	params := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"tacoIdentifier": {
				S: aws.String(tacoIdentifier),
			},
		},
		TableName: aws.String(h.Table),
	}
	_, err := h.Connection.DeleteItem(params)
	if err != nil {
		return err
	}
	return nil
}

// DestroyAll removes all rows from the repository
// We're only using this for testing in order to reset the
// database to a known clean state
func (h DynamodbDatabase) DestroyAll() error {
	projection := "tacoIdentifier"
	out, _ := h.Connection.Scan(&dynamodb.ScanInput{
		TableName:            aws.String(h.Table),
		ProjectionExpression: &projection,
	})
	for _, element := range out.Items {
		if err := h.DeleteByID(*element[projection].S); err != nil {
			panic(err)
		}
	}
	return nil
}
