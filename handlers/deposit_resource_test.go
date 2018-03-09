package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
)

func TestCreateResourceHappyPath(t *testing.T) {
	r := gofight.New()
	repo := mockRepo(nil)

	byt, err := ioutil.ReadFile("../examples/create-bs646cd8717.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	r.POST("/v1/resource").
		SetJSON(postData).
		Run(setupFakeRuntime().WithRepository(repo).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*fakeRepository).CreatedResources))
				// TODO: ResourceIdentification has no SourceID?
				// assert.Equal(t, "bib12345678", repo.(*fakeRepository).CreatedResources[0].SourceID)
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":    "oo000oo0001",
			"title": "My work",
		}).
		Run(setupFakeRuntime().Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

func TestCreateResourceSemanticallyValid(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"id":       "oo000oo0001",
			"@context": "http://example.com", // This is not a valid context
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"label":    "My work"}).
		Run(setupFakeRuntime().Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

func TestCreateResourceFailure(t *testing.T) {
	r := gofight.New()

	byt, err := ioutil.ReadFile("../examples/create-bs646cd8717.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	r.POST("/v1/resource").
		SetJSON(postData).
		Run(setupFakeRuntime().WithRepository(mockErrorRepo()).Handler(),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				log.Println(r.Body.String())
				assert.Equal(t, 500, r.Code)
			})
}
