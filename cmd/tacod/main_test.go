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
	"github.com/sul-dlss-labs/taco/persistence"
	"github.com/sul-dlss-labs/taco/streaming"
)

func mockRepo(record *persistence.Resource) persistence.Repository {
	return &fakeRepository{record: record}
}

type fakeRepository struct {
	record *persistence.Resource
}

func (f fakeRepository) GetByID(id string) (*persistence.Resource, error) {

	if f.record != nil {
		return f.record, nil
	}
	return nil, errors.New("not found")
}

func (f fakeRepository) SaveItem(resource *persistence.Resource) {}

func mockStream() streaming.Stream {
	return &fakeStream{}
}

type fakeStream struct {
}

func (d fakeStream) SendMessage(message string)                      {}
func (d fakeStream) GetIterator(shardID *string) *string             { return nil }
func (d fakeStream) GetRecords(iterator *string) ([]string, *string) { return nil, nil }

func setupFakeRuntime(repo persistence.Repository) http.Handler {
	config.Init("../../config/test.yaml")
	rt, _ := taco.NewRuntimeWithServices(viper.GetViper(), repo, mockStream())
	return buildAPI(rt).Serve(nil)
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

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":       "oo000oo0001",
			"sourceId": "bib12345678",
			"title":    "My work",
		}).
		Run(setupFakeRuntime(mockRepo(nil)),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(setupFakeRuntime(mockRepo(nil)),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}
