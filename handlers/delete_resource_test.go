package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestDeleteResourceHappyPath(t *testing.T) {
	r := gofight.New()
	json := datautils.JSONObject{}
	repo := NewMockDatabase(datautils.NewResource(json).WithID("99999"))
	r.DELETE("/v1/resource/oo000oo0001").
		Run(handler(repo, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNoContent, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).DeletedResources))
				assert.Equal(t, "oo000oo0001", repo.(*MockDatabase).DeletedResources[0])
			})
}

func TestDeleteResourceFailure(t *testing.T) {
	r := gofight.New()
	repo := NewMockErrorDatabase()
	assert.Panics(t,
		func() {
			r.DELETE("/v1/resource/oo000oo0001").
				Run(handler(repo, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
