package validators

import (
	"errors"
	"io/ioutil"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
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

func (f *fakeRepository) Read(id string) (*datautils.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) Insert(resource datautils.Resource) error {
	return errors.New("not implemented")
}

func (d *fakeRepository) UpdateString(resourceID string, field string, value string) error {
	return errors.New("not implemented")
}

func (d *fakeRepository) UpdateBool(resourceID string, field string, value bool) error {
	return errors.New("not implemented")
}
