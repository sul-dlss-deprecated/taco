package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

var updateMessage = gofight.D{
	"id":       "99",
	"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
	"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
	"access": gofight.D{
		"access":   "world",
		"download": "world",
	},
	"administrative": gofight.D{
		"sdrPreserve":   false,
		"isDescribedBy": "the_mods.xml",
		"created":       "2012-10-28T04:13:43.639Z",
	},
	"currentVersion": true,
	"version":        5,
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
	repo := NewMockDatabase(&datautils.Resource{"id": "99"})

	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(updateMessage).
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestUpdateResourceNotFound(t *testing.T) {
	r := gofight.New()
	r.PATCH("/v1/resource/99").
		SetHeader(gofight.H{"Content-Type": "application/json"}).
		SetJSON(updateMessage).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
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
