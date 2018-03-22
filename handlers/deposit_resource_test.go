package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func postData() map[string]interface{} {
	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}
	return postData
}

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(nil)
	stream := NewMockStream("")

	r.POST("/v1/resource").
		SetJSON(postData()).
		Run(handler(repo, stream, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				resource := repo.(*MockDatabase).CreatedResources[0]
				assert.Equal(t, 1, resource.Version())
				assert.True(t, resource.CurrentVersion())

				// assert.Equal(t, "bib12345678", repo.(*MockDatabase).CreatedResources[0].(map[string]interface{})["sourceid"])
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(handler(NewMockDatabase(nil), NewMockStream(""), nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

func TestCreateInvalidResource(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":       "oo000oo0001",
			"@context": "http://example.com", // This is not a valid context
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(handler(nil, nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
				assert.Contains(t, r.Body.String(), "Validation Error")
			})
}

func TestCreateResourceFailure(t *testing.T) {
	r := gofight.New()
	assert.Panics(t,
		func() {
			r.POST("/v1/resource").
				SetJSON(postData()).
				Run(handler(NewMockErrorDatabase(), NewMockStream(""), nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
