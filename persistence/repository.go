package persistence

import (
	"errors"
	"log"
	"reflect"

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
	CreateItem(*Resource) error
	UpdateItem(*Resource, *Resource) error
}

// DynamoRepository -- The fetching object
type DynamoRepository struct {
	db        *dynamodb.DynamoDB
	tableName *string
}

// SaveItem perist the resource in dynamo db
func (h DynamoRepository) CreateItem(resource *Resource) error {
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

// UpdateItem - Compares provided resource to existing resource and calls the appropriate update method
func (h DynamoRepository) UpdateItem(resource *Resource, originalResource *Resource) error {

	var err error

	res := reflect.ValueOf(resource).Elem()
	typeOfResource := res.Type()
	orig := reflect.ValueOf(originalResource).Elem()

	for i := 0; i < res.NumField(); i++ {
		resourceField := res.Field(i)
		origResourceField := orig.Field(i)

		fieldName := typeOfResource.Field(i).Name

		if resourceField.Interface() != origResourceField.Interface() {
			switch resourceField.Type().String() {
			case "bool":
				err = h.updateBool(resource.ID, fieldName, resourceField.Interface().(bool))
			default:
				err = h.updateString(resource.ID, fieldName, resourceField.Interface().(string))
			}
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// updateString - Updates a string value in the dynamo repository for this resource
func (h DynamoRepository) updateString(resourceID string, field string, value string) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(value),
			},
		},
		TableName: h.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(resourceID),
			},
		},
		UpdateExpression: aws.String("set " + field + " = :r"),
	}

	_, err := h.db.UpdateItem(input)

	if err != nil {
		return err
	}

	return nil
}

// updateBool - Updates a boolean value in the dynamo repository for this resource
func (h DynamoRepository) updateBool(resourceID string, field string, value bool) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				BOOL: aws.Bool(value),
			},
		},
		TableName: h.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(resourceID),
			},
		},
		UpdateExpression: aws.String("set " + field + " = :r"),
	}

	_, err := h.db.UpdateItem(input)

	if err != nil {
		return err
	}

	return nil
}
