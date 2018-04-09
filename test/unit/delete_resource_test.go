package handlers

import (
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func TestDeleteResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(nil)
	stream := NewMockStream()
	r.DELETE("/v1/resource/oo000oo0001").
		Run(handler(repo, stream, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusNoContent, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).DeletedResources))
				assert.Equal(t, "oo000oo0001", repo.(*MockDatabase).DeletedResources[0])
				assert.Equal(t, 1, len(stream.(*MockStream).Messages))
				assert.Equal(t, "delete", stream.(*MockStream).Messages[0].Action)
				assert.Equal(t, "oo000oo0001", stream.(*MockStream).Messages[0].ID)
			})
}

func TestDeleteResourceFailure(t *testing.T) {
	r := gofight.New()
	repo := NewMockErrorDatabase()
	assert.Panics(t,
		func() {
			r.DELETE("/v1/resource/oo000oo0001").
				Run(handler(repo, nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
