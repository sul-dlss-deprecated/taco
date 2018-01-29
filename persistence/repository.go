package persistence

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/config"
)

// NewDynamoRepository -- Creates a new repository
func NewDynamoRepository(config *config.Config, db *dynamodb.DynamoDB) *DynamoRepository {
	tableName := aws.String(config.ResourceTableName)
	return &DynamoRepository{db: db,
		tableName: tableName}
}

// Repository the interface for the metadata repository
type Repository interface {
	GetByID(string) (*Resource, error)
	SaveItem(*Resource) error
}

// DynamoRepository -- The fetching object
type DynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName *string
}

// SaveItem perist the resource in dynamo db
func (h DynamoRepository) SaveItem(resource *Resource) error {
	row, err := dynamodbattribute.MarshalMap(resource)

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: h.tableName,
	}

	_, err = h.db.PutItem(input)

	if err != nil {
		return err
	}
	log.Printf("Saved %s to dynamodb", resource.ID)
	return nil
}

// GetByID -- given an identifier, find the resource
func (h DynamoRepository) GetByID(id string) (*Resource, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:      h.tableName,
		ConsistentRead: aws.Bool(true),
	}
	resp, err := h.db.GetItem(params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var resource *Resource
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &resource); err != nil {
		log.Println(err)
		return nil, err
	}

	if resource.ID == "" {
		return nil, errors.New("not found")
	}
	return resource, nil
}
