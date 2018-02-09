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
	repo := NewMockDatabase(datautils.NewResource(json))
	r.GET("/v1/file/99").
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusFound, r.Code)
				assert.Equal(t, "https://example.com/file-123", r.HeaderMap.Get("Location"))
			})
}

func TestRetrieveFileNotFound(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/file/100").
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}

func TestRetrieveFileError(t *testing.T) {
	r := gofight.New()
	repo := NewMockErrorDatabase()

	assert.Panics(t,
		func() {
			r.GET("/v1/file/100").
				Run(handler(repo, nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
						assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
					})
		})
}
