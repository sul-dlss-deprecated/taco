package handlers

import (
	"errors"

	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

type MockDatabase struct {
	CreatedResources []interface{}
}

func NewMockDatabase() db.Database {
	return &MockDatabase{CreatedResources: []interface{}{}}
}

func (d *MockDatabase) Insert(params interface{}) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) Read(id string) (*models.Resource, error) {
	return nil, nil
}

type MockErrorDatabase struct {
}

func NewMockErrorDatabase() db.Database {
	return &MockDatabase{}
}

func (d *MockErrorDatabase) Insert(params interface{}) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) Read(id string) (*models.Resource, error) {
	return nil, errors.New("Broken")
}
