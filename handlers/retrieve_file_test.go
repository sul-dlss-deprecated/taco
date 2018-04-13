package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestRetrieveFileHappyPath(t *testing.T) {
	r := gofight.New()
	json := datautils.JSONObject{"file-location": "s3://bucket/key"}
	repo := NewMockDatabase(datautils.NewResource(json), nil)
	r.GET("/v1/file/99").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusFound, r.Code)
				assert.Equal(t, "https://example.com/file-123", r.HeaderMap.Get("Location"))
			})
}

func TestRetrieveFileNotFound(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/file/100").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}

func TestRetrieveFileNoBehalfOfHeader(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/file/99").
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestRetrieveFileNoPermissions(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(new(datautils.Resource), nil)
	r.GET("/v1/file/99").
		SetHeader(gofight.H{
			"On-Behalf-Of": "blalbrit@stanford.edu", // The dummy authZ service is set to only allow lmcrae@stanford.edu
		}).
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestRetrieveFileError(t *testing.T) {
	r := gofight.New()
	repo := NewMockErrorDatabase(nil)

	assert.Panics(t,
		func() {
			r.GET("/v1/file/100").
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
				}).
				Run(handler(repo, nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
						assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
					})
		})
}
