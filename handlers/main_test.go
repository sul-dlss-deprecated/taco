package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/persistence"
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

func setupFakeRuntime(repo persistence.Repository) http.Handler {
	config.Init("../config/test.yaml")
	rt, _ := taco.NewRuntimeForRepository(viper.GetViper(), repo)
	return BuildAPI(rt).Serve(nil)
}

func TestRetrieveHappyPath(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/resource/99").
		Run(setupFakeRuntime(mockRepo(new(persistence.Resource))),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestRetrieveNotFound(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/resource/100").
		Run(setupFakeRuntime(mockRepo(nil)),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}
