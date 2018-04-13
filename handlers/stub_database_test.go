package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
)

func handler(database db.Database, storage storage.Storage, identifierService identifier.Service) http.Handler {
	if database == nil {
		database = NewMockDatabase(nil, nil)
	}
	if storage == nil {
		storage = NewMockStorage()
	}

	if identifierService == nil {
		identifierService = identifier.NewUUIDService()
	}

	return BuildAPI(database, storage, identifierService).Serve(nil)
}

type MockDatabase struct {
	current          *datautils.Resource
	specificVersion  *datautils.Resource
	CreatedResources []*datautils.Resource
	DeletedResources []string
}

func NewMockDatabase(current *datautils.Resource, specificVersion *datautils.Resource) db.Database {
	return &MockDatabase{
		CreatedResources: []*datautils.Resource{},
		DeletedResources: []string{},
		current:          current,
		specificVersion:  specificVersion,
	}
}

func (d *MockDatabase) Insert(params *datautils.Resource) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) RetrieveLatest(externalID string) (*datautils.Resource, error) {
	if d.current != nil {
		record := d.current
		d.current = nil
		return record, nil
	}
	return nil, &db.RecordNotFound{ID: &externalID}
}

func (d *MockDatabase) RetrieveVersion(externalID string, version *string) (*datautils.Resource, error) {
	if d.specificVersion != nil {
		record := d.specificVersion
		d.specificVersion = nil
		return record, nil
	}
	return nil, errors.New("not found")
}

func (d *MockDatabase) Update(params *datautils.Resource) error {
	return nil
}

func (d *MockDatabase) DeleteVersion(tacoIdentifier string) error {
	d.DeletedResources = append(d.DeletedResources, tacoIdentifier)
	return nil
}

type MockErrorDatabase struct {
	record *datautils.Resource
}

func NewMockErrorDatabase(record *datautils.Resource) db.Database {
	return &MockErrorDatabase{
		record: record,
	}
}

func (d *MockErrorDatabase) Insert(params *datautils.Resource) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) Update(params *datautils.Resource) error {
	return nil
}

func (d *MockErrorDatabase) RetrieveLatest(externalID string) (*datautils.Resource, error) {
	return d.record, nil
}

func (d *MockErrorDatabase) DeleteByID(tacoIdentifier string) error {
	return errors.New("Broken")
}
func (d *MockErrorDatabase) DeleteVersion(tacoIdentifier string) error {
	return errors.New("Broken")
}

func (d *MockErrorDatabase) RetrieveVersion(id string, version *string) (*datautils.Resource, error) {
	return nil, errors.New("not implemented")
}
