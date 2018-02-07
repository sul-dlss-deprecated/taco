package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/persistence"
)

func TestUpdateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := mockRepo(new(persistence.Resource))

	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(gofight.D{
			"id":       "99",
			"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My updated work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(setupFakeRuntime().WithRepository(repo).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestUpdateResourceNotFound(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(gofight.D{
			"id":       "99",
			"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My updated work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(setupFakeRuntime().WithRepository(mockRepo(nil)).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}

func TestUpdateResourceEmptyRequest(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/100").
		Run(setupFakeRuntime().WithRepository(mockRepo(nil)).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}
