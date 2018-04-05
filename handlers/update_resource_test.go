package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

var updateMessageSameVersion = gofight.D{
	"externalIdentifier": "99",
	"@context":           "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
	"@type":              "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
	"tacoIdentifier":     "99",
	"access": gofight.D{
		"access":   "world",
		"download": "world",
	},
	"administrative": gofight.D{
		"sdrPreserve":   false,
		"isDescribedBy": "the_mods.xml",
		"created":       "2012-10-28T04:13:43.639Z",
	},
	"version": 1,
	"depositor": gofight.D{
		"name":    "Lynn",
		"sunetID": "lmcray",
	},
	"identification": gofight.D{
		"identifier": "1234abc",
		"sdrUUID":    "123888019239",
	},
	"structural": gofight.D{
		"hasAgreement": "yes",
	},
	"label": "My updated work"}

var updateMessageNewVersion = gofight.D{
	"externalIdentifier": "99",
	"@context":           "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
	"@type":              "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
	"tacoIdentifier":     "99",
	"access": gofight.D{
		"access":   "world",
		"download": "world",
	},
	"administrative": gofight.D{
		"sdrPreserve":   false,
		"isDescribedBy": "the_mods.xml",
		"created":       "2012-10-28T04:13:43.639Z",
	},
	"version": 2,
	"depositor": gofight.D{
		"name":    "Lynn",
		"sunetID": "lmcray",
	},
	"identification": gofight.D{
		"identifier": "1234abc",
		"sdrUUID":    "123888019239",
	},
	"structural": gofight.D{
		"hasAgreement": "yes",
	},
	"label": "My updated work"}

func TestUpdateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(datautils.NewResource(datautils.JSONObject{"tacoIdentifier": "99", "version": float64(1), "externalIdentifier": "99"}))

	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(updateMessageSameVersion).
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				resource := repo.(*MockDatabase).CreatedResources[0]
				assert.Equal(t, 1, resource.Version())
			})
}

func TestUpdateResourceNewVersionPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(datautils.NewResource(datautils.JSONObject{"tacoIdentifier": "99", "version": float64(1), "externalIdentifier": "99"}))

	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(updateMessageNewVersion).
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, 2, len(repo.(*MockDatabase).CreatedResources))
				oldResource := repo.(*MockDatabase).CreatedResources[1]
				newResource := repo.(*MockDatabase).CreatedResources[0]
				assert.Equal(t, 1, oldResource.Version())
				assert.Equal(t, newResource.ID(), oldResource.JSON["followingVersion"])
				assert.Equal(t, 2, newResource.Version())
				assert.Equal(t, "99", newResource.JSON["precedingVersion"])
			})
}

func TestUpdateResourceNotFound(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(updateMessageSameVersion).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}

func TestUpdateInvalidResource(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/100").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(gofight.D{
			"externalIdentifier": "oo000oo0001",
			"@context":           "http://example.com", // This is not a valid context
			"@type":              "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":             "world",
			"label":              "My work",
			"preserve":           true,
			"publish":            true,
			"sourceId":           "bib12345678"}).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
				assert.Contains(t, r.Body.String(), "Validation Error")
			})
}

func TestUpdateResourceEmptyRequest(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/100").
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}
