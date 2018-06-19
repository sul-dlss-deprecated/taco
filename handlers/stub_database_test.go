package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/identifier"
	"github.com/sul-dlss-labs/taco/storage"
	"github.com/sul-dlss-labs/taco/validators"
)

type runtime struct {
	database          db.Database
	storage           storage.Storage
	identifierService identifier.Service
	updateValidator   validators.ResourceValidator
	depositValidator  validators.ResourceValidator
	fileValidator     validators.ResourceValidator
}

type dummySuccessValidator struct {
}

func (v *dummySuccessValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	return nil
}

type dummyFailValidator struct {
}

func (v *dummyFailValidator) ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors {
	errors := models.ErrorResponseErrors{}
	source := &models.ErrorSource{Pointer: "structural.hasMember"}
	problem := &models.Error{
		Title:  "Validation Error",
		Detail: "Stub error",
		Source: source,
	}
	errors = append(errors, problem)
	return &errors
}

func handler(opts ...*runtime) http.Handler {
	var rt *runtime
	if len(opts) == 0 || opts[0] == nil {
		rt = &runtime{}
	} else {
		rt = opts[0]
	}
	if rt.database == nil {
		rt.database = NewMockDatabase(nil)
	}
	if rt.storage == nil {
		rt.storage = NewMockStorage()
	}

	if rt.identifierService == nil {
		rt.identifierService = identifier.NewUUIDService()
	}

	if rt.depositValidator == nil {
		rt.depositValidator = &dummySuccessValidator{}
	}

	if rt.updateValidator == nil {
		rt.updateValidator = &dummySuccessValidator{}
	}

	if rt.fileValidator == nil {
		rt.fileValidator = &dummySuccessValidator{}
	}

	return BuildAPI(rt.database, rt.storage, rt.identifierService, rt.depositValidator, rt.updateValidator, rt.fileValidator).Serve(nil)
}

type MockDatabase struct {
	record           *datautils.Resource
	CreatedResources []*datautils.Resource
	DeletedResources []string
}

func NewMockDatabase(record *datautils.Resource) db.Database {
	return &MockDatabase{
		CreatedResources: []*datautils.Resource{},
		DeletedResources: []string{},
		record:           record,
	}
}

func (d *MockDatabase) Insert(params *datautils.Resource) error {
	d.CreatedResources = append(d.CreatedResources, params)
	return nil
}

func (d *MockDatabase) RetrieveLatest(externalID string) (*datautils.Resource, error) {
	if d.record != nil {
		record := d.record
		d.record = nil
		return record, nil
	}
	return nil, &db.RecordNotFound{ID: &externalID}
}

func (d *MockDatabase) RetrieveVersion(externalID string, version *string) (*datautils.Resource, error) {
	return nil, errors.New("not implemented")
}

func (d *MockDatabase) Update(params *datautils.Resource) error {
	return nil
}

func (d *MockDatabase) DeleteByID(tacoIdentifier string) error {
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

func (d *MockErrorDatabase) RetrieveVersion(id string, version *string) (*datautils.Resource, error) {
	return nil, errors.New("not implemented")
}
