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
	resource := datautils.NewResource(json).
		WithID("99999").
		WithType(datautils.ObjectTypes[0])
	repo := NewMockDatabase(resource)
	r.DELETE("/v1/resource/oo000oo0001").
		Run(handler(&runtime{database: repo}),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNoContent, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).DeletedResources))
				assert.Equal(t, "99999", repo.(*MockDatabase).DeletedResources[0])
			})
}

func TestDeleteFileResourceHappyPath(t *testing.T) {
	r := gofight.New()
	json := datautils.JSONObject{}
	resource := datautils.NewResource(json).
		WithID("99999").
		WithType(datautils.FileType).
		WithFileLocation("s3://bucket/my-key")
	repo := NewMockDatabase(resource)
	storage := NewMockStorage()

	r.DELETE("/v1/resource/oo000oo0001").
		Run(handler(repo, storage, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNoContent, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).DeletedResources))
				assert.Equal(t, "99999", repo.(*MockDatabase).DeletedResources[0])
				assert.Equal(t, "s3://bucket/my-key", storage.(*MockStorage).DeletedFiles[0])
			})
}

func TestDeleteResourceFailure(t *testing.T) {
	r := gofight.New()
	resource := datautils.NewResource(datautils.JSONObject{}).
		WithID("99999")
	repo := NewMockErrorDatabase(resource)
	assert.Panics(t,
		func() {
			r.DELETE("/v1/resource/oo000oo0001").
				Run(handler(&runtime{database: repo}),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
