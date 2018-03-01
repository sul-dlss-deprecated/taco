package persistence

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// PrimaryKey is the primary key of the table
const PrimaryKey = "id"

// NewDynamoRepository -- Creates a new repository. This is an implentation of
// an abstract "repository" concept that is backed by a single table
// (config.ResourceTableName) in DynamoDB.
func NewDynamoRepository(config *config.Config, db *dynamodb.DynamoDB) *DynamoRepository {
	tableName := aws.String(config.ResourceTableName)
	return &DynamoRepository{db: db,
		tableName: tableName}
}

// Repository the interface for the metadata repository
type Repository interface {
	GetByID(string) (*models.Resource, error)
	CreateItem(map[string]*dynamodb.AttributeValue) error
	UpdateItem(map[string]*dynamodb.AttributeValue) error
}

// DynamoRepository -- The fetching object
type DynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName *string
}

// CreateItem perist the resource in dynamo db
func (h DynamoRepository) CreateItem(row map[string]*dynamodb.AttributeValue) error {

	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: h.tableName,
	}

	_, err := h.db.PutItem(input)

	if err != nil {
		return err
	}
	log.Printf("Saved %s to dynamodb", row[PrimaryKey])
	return nil
}

// GetByID -- given an identifier, find the resource
func (h DynamoRepository) GetByID(id string) (*models.Resource, error) {
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

// UpdateItem - Replaces an existing resource in the repository
func (h DynamoRepository) UpdateItem(row map[string]*dynamodb.AttributeValue) error {
	input := &dynamodb.PutItemInput{
		Item:      row,
		TableName: h.tableName,
	}

	_, err := h.db.PutItem(input)

	if err != nil {
		return err
	}
	return nil
}
