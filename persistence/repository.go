package persistence

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/config"
)

// NewRepository -- Creates a new repository
func NewRepository(db *dynamodb.DynamoDB) (*DynamoRepository, error) {
	config := config.NewConfig()
	tableName := aws.String(config.Resource_Table_Name)
	return &DynamoRepository{db: db,
			tableName: tableName},
		nil
}

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
