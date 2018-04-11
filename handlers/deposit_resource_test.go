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

	r.POST("/v1/resource").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		SetJSON(postData()).
		Run(handler(repo, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				resource := repo.(*MockDatabase).CreatedResources[0]
				assert.Equal(t, 1, resource.Version())

				// assert.Equal(t, "bib12345678", repo.(*MockDatabase).CreatedResources[0].(map[string]interface{})["sourceid"])
			})
}

func TestCreateResourceNoApiKey(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"tacoIdentifier":       "oo000oo0001",
			"sourceId": "bib12345678",
			"title":    "My work",
		}).
		Run(handler(nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateResourceNoPermissions(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetHeader(gofight.H{
			"On-Behalf-Of": "blalbrit@stanford.edu", // The dummy authZ service is set to only allow lmcrae@stanford.edu
		}).
		SetJSON(postData()).
		Run(handler(nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnauthorized, r.Code)
			})
}

func TestCreateResourceMissingSourceId(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		SetJSON(gofight.D{
			"tacoIdentifier":    "oo000oo0001",
			"@type": "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"title": "My work",
		}).
		Run(handler(nil, nil),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusUnprocessableEntity, r.Code)
			})
}

func TestCreateInvalidResource(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		SetJSON(gofight.D{
			"tacoIdentifier":       "oo000oo0001",
			"@context": "http://example.com", // This is not a valid context
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"label":    "My work",
			"preserve": true,
			"publish":  true,
			"sourceId": "bib12345678"}).
		Run(handler(nil, nil),
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
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
				}).
				SetJSON(postData()).
				Run(handler(NewMockErrorDatabase(), nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
