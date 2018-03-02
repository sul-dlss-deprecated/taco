package validators

import (
	"errors"

	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/serializers"
)

func testResource() *models.Resource {
	return &models.Resource{}
}

func newMockRepository() persistence.Repository {
	return &fakeRepository{}
}

type fakeRepository struct{}

func (f *fakeRepository) GetByID(id string) (*models.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) CreateItem(resource *serializers.Resource) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) UpdateItem(resource *serializers.Resource) error {
	return errors.New("not implemented")
}
