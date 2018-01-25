package persistence

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/spf13/viper"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// NewRepository -- Creates a new repository
func NewRepository(config viper.Viper, db *dynamodb.DynamoDB) (*DynamoRepository, error) {
	return &DynamoRepository{config: config, db: db}, nil
}

type Repository interface {
	GetByID(string) (*models.Resource, error)
}

// DynamoRepository -- The fetching object
type DynamoRepository struct {
	config viper.Viper
	db     *dynamodb.DynamoDB
}

// GetByID -- given an identifier, find the resource
func (h DynamoRepository) GetByID(id string) (*models.Resource, error) {
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
	var resource *models.Resource
	if err := dynamodbattribute.UnmarshalMap(resp.Item, &resource); err != nil {
		log.Println(err)
		return nil, err
	}

	if resource.ID == "" {
		return nil, errors.New("not found")
	}
	return resource, nil
}
