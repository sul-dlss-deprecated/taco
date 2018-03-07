package validators

import (
	"errors"
	"io/ioutil"

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

func newMockRepository() persistence.Repository {
	return &fakeRepository{}
}

type fakeRepository struct{}

func (f *fakeRepository) GetByID(id string) (*persistence.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) CreateItem(resource persistence.Resource) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) UpdateItem(resource persistence.Resource) error {
	return errors.New("not implemented")
}
