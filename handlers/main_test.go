package handlers

import (
	"errors"
	"net/http"

	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
)

func mockRepo(record *persistence.Resource) persistence.Repository {
	return &fakeRepository{record: record, CreatedResources: []persistence.Resource{}}
}

type fakeRepository struct {
	record           *persistence.Resource
	CreatedResources []persistence.Resource
}

func (f *fakeRepository) GetByID(id string) (*persistence.Resource, error) {

	if f.record != nil {
		return f.record, nil
	}
	return nil, errors.New("not found")
}

func (f *fakeRepository) SaveItem(resource *persistence.Resource) error {
	f.CreatedResources = append(f.CreatedResources, *resource)
	return nil
}

func mockStream() streaming.Stream {
	return &fakeStream{}
}

type fakeStream struct {
}

func (d fakeStream) SendMessage(message string) error { return nil }

func setupFakeRuntime(repo persistence.Repository) http.Handler {
	rt, _ := taco.NewRuntimeWithServices(config.NewConfig(), repo, mockStream())
	return BuildAPI(rt).Serve(nil)
}

func mockErrorRepo() persistence.Repository {
	return &fakeErroringRepository{}
}

type fakeErroringRepository struct{}

func (f *fakeErroringRepository) GetByID(id string) (*persistence.Resource, error) {
	return nil, errors.New("broken")
}

func (f *fakeErroringRepository) SaveItem(resource *persistence.Resource) error {
	return errors.New("broken")
}
