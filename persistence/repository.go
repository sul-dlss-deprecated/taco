package persistence

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/spf13/viper"
)

// NewRepository -- Creates a new repository
func NewRepository(config viper.Viper, db *dynamodb.DynamoDB) (*DynamoRepository, error) {
	return &DynamoRepository{config: config, db: db}, nil
}

// Repository -- an interface for the metadata repository
type Repository interface {
	GetByID(string) (*Resource, error)
	SaveItem(*Resource)
}

// DynamoRepository -- The fetching object
type DynamoRepository struct {
	config viper.Viper
	db     *dynamodb.DynamoDB
}

// SaveItem perist the resource in dynamo db
func (h DynamoRepository) SaveItem(resource *Resource) {
}

// GetByID -- given an identifier, find the resource
func (h DynamoRepository) GetByID(id string) (*Resource, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:      aws.String("resources"),
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
