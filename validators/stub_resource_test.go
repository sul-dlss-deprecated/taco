package validators

import (
	"errors"
	"io/ioutil"

	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/persistence"
)

func testResource() string {
	byt, err := ioutil.ReadFile("../examples/bs646cd8717.json")
	if err != nil {
		panic(err)
	}
	return string(byt)
}

func testDepositResource() string {
	byt, err := ioutil.ReadFile("../examples/create-bs646cd8717.json")
	if err != nil {
		panic(err)
	}
	return string(byt)
}

func newMockRepository() db.Database {
	return &fakeRepository{}
}

type fakeRepository struct{}

func (f *fakeRepository) Read(id string) (*models.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) Insert(resource persistence.Resource) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) Update(resource persistence.Resource) error {
	return errors.New("not implemented")
}
