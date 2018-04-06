package validators

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/db"
)

func testResource(file string) *datautils.Resource {
	byt, err := ioutil.ReadFile(fmt.Sprintf("../examples/%s", file))
	if err != nil {
		panic(err)
	}
	data := datautils.JSONObject{}
	if err := json.Unmarshal(byt, &data); err != nil {
		panic(err)
	}

	return datautils.NewResource(data)
}

func newMockRepository(record *datautils.Resource) db.Database {
	return &fakeRepository{
		record: record,
	}
}

type fakeRepository struct {
	record *datautils.Resource
}

func (f *fakeRepository) RetrieveLatest(externalID string) (*datautils.Resource, error) {
	if f.record != nil {
		return f.record, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepository) RetrieveVersion(externalID string, version *string) (*datautils.Resource, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeRepository) Insert(resource *datautils.Resource) error {
	return errors.New("not implemented")
}

func (f *fakeRepository) DeleteAllVersions(externalID string) error {
	return errors.New("not implemented")
}
