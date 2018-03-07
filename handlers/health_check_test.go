package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/buger/jsonparser"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	r := gofight.New()
	r.GET("/v1/healthcheck").
		Run(setupFakeRuntime().Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusOK, r.Code)
				stat, _ := jsonparser.GetString(r.Body.Bytes(), "status")
				assert.Equal(t, "OK", stat)
			})
}
