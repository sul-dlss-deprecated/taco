package validators

import (
	"errors"

	"github.com/sul-dlss-labs/taco/persistence"
)

func testResource() string {
	return ""
}

func testDepositResource() string {
	return ""
}

func newMockRepository() persistence.Repository {
	return &fakeRepository{}
}

type fakeRepository struct{}

func (f *fakeRepository) GetByID(id string) (*persistence.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) CreateItem(resource *persistence.Resource) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) UpdateItem(resource *persistence.Resource) error {
	return errors.New("not implemented")
}
