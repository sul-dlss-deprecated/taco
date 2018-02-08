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
		Run(setupFakeRuntime().WithRepository(repo).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*fakeRepository).CreatedResources))
				assert.Equal(t, "bib12345678", repo.(*fakeRepository).CreatedResources[0].SourceID)
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(setupFakeRuntime().Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

// TODO: Handle errors
func TestCreateResourceFailure(t *testing.T) {
	r := gofight.New()
	assert.Panics(t,
		func() {
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
				Run(setupFakeRuntime().WithRepository(mockErrorRepo()).Handler(),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
