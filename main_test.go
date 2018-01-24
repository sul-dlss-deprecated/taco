package main

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/server"
)

func TestRetrieveHappyPath(t *testing.T) {
	config.Init("test")
	db.Init()

	r := gofight.New()

	r.GET("/v1/resource/99").
		SetDebug(true).
		Run(server.NewRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, http.StatusOK, r.Code)
		})
}

func TestRetrieveNotFound(t *testing.T) {
	config.Init("test")
	db.Init()

	r := gofight.New()

	r.GET("/v1/resource/100").
		SetDebug(true).
		Run(server.NewRouter(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, 404, r.Code)
		})
}
