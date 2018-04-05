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
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		Run(handler(repo, nil, nil),
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
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
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
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
				}).
				Run(handler(repo, nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}

func TestDeleteResourceNoApiKey(t *testing.T) {
	r := gofight.New()
	r.DELETE("/v1/resource/oo000oo0001").
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestDeleteResourceNoPermissions(t *testing.T) {
	r := gofight.New()
	r.DELETE("/v1/resource/oo000oo0001").
		SetHeader(gofight.H{
			"On-Behalf-Of": "blalbrit@stanford.edu", // The dummy authZ service is set to only allow lmcrae@stanford.edu
		}).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}
