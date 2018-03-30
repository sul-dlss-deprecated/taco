package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestRetrieveHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(&datautils.Resource{})
	r.GET("/v1/resource/99").
		Run(handler(repo, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
			})
}

func TestRetrieveNotFound(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/resource/100").
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNotFound, r.Code)
			})
}

func TestRetrieveError(t *testing.T) {
	r := gofight.New()
	assert.Panics(t,
		func() {
			r.GET("/v1/resource/100").
				Run(handler(NewMockErrorDatabase(), nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})

		})
}
