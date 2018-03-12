package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

func handler(database db.Database) http.Handler {
	if database == nil {
		database = NewMockDatabase(nil)
	}
	return BuildAPI(database).Serve(nil)
}

type MockDatabase struct {
	record           *models.Resource
	CreatedResources []interface{}
}

func NewMockDatabase(record *models.Resource) db.Database {
	return &MockDatabase{CreatedResources: []interface{}{}, record: record}
}

func (d *MockDatabase) Insert(params interface{}) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) Read(id string) (*models.Resource, error) {
	if d.record != nil {
		return d.record, nil
	}
	return nil, errors.New("not found")
}

func (d *MockDatabase) Update(params interface{}) error {
	return nil
}

type MockErrorDatabase struct {
}

func NewMockErrorDatabase() db.Database {
	return &MockErrorDatabase{}
}

func (d *MockErrorDatabase) Insert(params interface{}) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) Update(params interface{}) error {
	return nil
}

func (d *MockErrorDatabase) Read(id string) (*models.Resource, error) {
	return nil, errors.New("Broken")
}
