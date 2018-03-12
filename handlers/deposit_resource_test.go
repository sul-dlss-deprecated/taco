package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/db"
)

func handler(database db.Database) http.Handler {
	if database == nil {
		database = mockRepo()
	}
	return BuildAPI(database).Serve(nil)
}

func mockRepo() db.Database {
	return &MockDatabase{}
}

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := mockRepo()

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
		Run(handler(repo),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				assert.Equal(t, "bib12345678", repo.(*MockDatabase).CreatedResources[0].(map[string]interface{})["sourceid"])
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(handler(mockRepo()),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

func TestCreateResourceSemanticallyValid(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":       "oo000oo0001",
			"@context": "http://example.com", // This is not a valid context
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(handler(mockRepo()),
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
				Run(handler(NewMockErrorDatabase()),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
