package validators

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
)

func testResource(file string) string {
	byt, err := ioutil.ReadFile(fmt.Sprintf("../examples/%s", file))
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

func (f *fakeRepository) Update(resource datautils.Resource) error {
	return errors.New("not implemented")
}
