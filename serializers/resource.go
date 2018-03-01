package serializers

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/persistence"
)

// ToPersistable converts the model Resource into a persisable object
func ToPersistable(resourceID string, request *models.Resource) map[string]*dynamodb.AttributeValue {
	row, err := dynamodbattribute.MarshalMap(request)
	if err != nil {
		panic(err)
	}
	// add the identifier
	row[persistence.PrimaryKey] = &dynamodb.AttributeValue{S: &resourceID}

	return row
}
