package main

import (
	"errors"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/generated/models"
	"github.com/sul-dlss-labs/taco/persistence"
)

func mockRepo(record *models.Resource) persistence.Repository {
	return &fakeRepository{record: record}
}

type fakeRepository struct {
	record *models.Resource
}

func (f fakeRepository) GetByID(id string) (*models.Resource, error) {

	if f.record != nil {
		return f.record, nil
	}
	return nil, errors.New("not found")
}

func setupFakeRuntime(repo persistence.Repository) http.Handler {
	config.Init("../../config/test.yaml")
	rt, _ := taco.NewRuntimeForRepository(viper.GetViper(), repo)
	return buildAPI(rt).Serve(nil)
}

func TestRetrieveHappyPath(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/resource/99").
		Run(setupFakeRuntime(mockRepo(new(models.Resource))),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestRetrieveNotFound(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/resource/100").
		Run(setupFakeRuntime(mockRepo(nil)),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, 404, r.Code)
			})
}
