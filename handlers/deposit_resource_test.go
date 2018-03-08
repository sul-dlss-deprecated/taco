package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(nil)
	stream := NewMockStream("")

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
		Run(handler(repo, stream, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
<<<<<<< HEAD
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				assert.Equal(t, "bib12345678", repo.(*MockDatabase).CreatedResources[0].(map[string]interface{})["sourceid"])
=======
				assert.Equal(t, 1, len(repo.(*fakeRepository).CreatedResources))
				// assert.Equal(t, "bib12345678", repo.(*fakeRepository).CreatedResources[0].SourceID)
>>>>>>> Validate using json schema
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(handler(NewMockDatabase(nil), NewMockStream(""), nil),
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
		Run(handler(NewMockDatabase(nil), NewMockStream(""), nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

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
				Run(handler(NewMockErrorDatabase(), NewMockStream(""), nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
