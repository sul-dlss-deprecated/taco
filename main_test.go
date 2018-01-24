package main

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveHappyPath(t *testing.T) {
	handler := buildAPI().Serve(nil)
	r := gofight.New()
	r.GET("/v1/resource/99").
		Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

// func TestRetrieveNotFound(t *testing.T) {
// 	handler := buildAPI().Serve(nil)
// 	r := gofight.New()
// 	r.GET("/v1/resource/100").
// 		SetDebug(true).
// 		Run(handler, func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
// 			assert.Equal(t, 404, r.Code)
// 		})
// }
