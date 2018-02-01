package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

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
