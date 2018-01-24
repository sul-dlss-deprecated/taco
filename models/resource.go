package models

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sul-dlss-labs/taco/db"
)

// Resource -- The metadata object
type Resource struct {
	ID     string `json:"id,omitempty"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// GetByID -- given an identifier, find the resource
func (h Resource) GetByID(id string) (*Resource, error) {
	db := db.GetDB()
	params := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName:      aws.String("resources"),
		ConsistentRead: aws.Bool(true),
	}
	log.Printf("Parms %s", params)
	resp, err := db.GetItem(params)
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
		log.Printf("Could not find item with id %s", id)
		return nil, errors.New("not found")
	}
	return resource, nil
}
