package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/streaming"
)

func handler(database db.Database, stream streaming.Stream, storage storage.Storage) http.Handler {
	if database == nil {
		database = NewMockDatabase(nil)
	}
	if stream == nil {
		stream = NewMockStream("")
	}
	if storage == nil {
		storage = NewMockStorage()
	}

	identifierService := identifier.NewUUIDService()
	return BuildAPI(database, stream, storage, identifierService).Serve(nil)
}

type MockDatabase struct {
	record           *datautils.Resource
	CreatedResources []*datautils.Resource
}

func NewMockDatabase(record *datautils.Resource) db.Database {
	return &MockDatabase{CreatedResources: []*datautils.Resource{}, record: record}
}

func (d *MockDatabase) Insert(params *datautils.Resource) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) Read(id string) (*datautils.Resource, error) {
	if d.record != nil {
		return d.record, nil
	}
	return nil, errors.New("not found")
}

func (d *MockDatabase) Update(params *datautils.Resource) error {
	return nil
}

type MockErrorDatabase struct {
}

func NewMockErrorDatabase() db.Database {
	return &MockErrorDatabase{}
}

func (d *MockErrorDatabase) Insert(params *datautils.Resource) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) Update(params *datautils.Resource) error {
	return nil
}

func (d *MockErrorDatabase) Read(id string) (*datautils.Resource, error) {
	return nil, errors.New("Broken")
}
