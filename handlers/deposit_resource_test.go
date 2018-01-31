package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/persistence"
)

func mockErrorRepo() persistence.Repository {
	return &fakeErroringRepository{}
}

type fakeErroringRepository struct{}

func (f *fakeErroringRepository) GetByID(id string) (*persistence.Resource, error) {
	return nil, nil
}

func (f *fakeErroringRepository) SaveItem(resource *persistence.Resource) error {
	return errors.New("not found")
}

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := mockRepo(nil)

	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":       "oo000oo0001",
			"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(setupFakeRuntime(repo),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*fakeRepository).CreatedResources))

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

// TODO: Handle errors
// func TestCreateResourceFailure(t *testing.T) {
// 	r := gofight.New()
// 	r.POST("/v1/resource").
// 		SetJSON(gofight.D{
// 			"id":       "oo000oo0001",
// 			"sourceId": "bib12345678",
// 			"title":    "My work",
// 		}).
// 		Run(setupFakeRuntime(mockErrorRepo()),
// 			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
// 			})
// }
