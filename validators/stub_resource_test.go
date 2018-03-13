package validators

import (
	"errors"

	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
)

func testResource() *models.Resource {
	return &models.Resource{}
}

func newMockRepository() db.Database {
	return &fakeRepository{}
}

type fakeRepository struct{}

func (f *fakeRepository) Read(id string) (*models.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) Insert(interface{}) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) Update(interface{}) error {
	return errors.New("not implemented")
}
