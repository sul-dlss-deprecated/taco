package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/appleboy/gofight"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

type FakeDruidService struct{}

func (d *FakeDruidService) Mint(*datautils.Resource) (string, error) { return "zt570tx3016", nil }

func postData(path string) map[string]interface{} {
	byt, err := ioutil.ReadFile(fmt.Sprintf("../examples/%s", path))
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}
	return postData
}

func TestCreateCollectionHappyPath(t *testing.T) {
	r := gofight.New()
	repo := NewMockDatabase(nil)
	idService := &FakeDruidService{}

	r.POST("/v1/resource").
		SetHeader(gofight.H{
			"On-Behalf-Of": "lmcrae@stanford.edu",
		}).
		SetJSON(postData("request.json")).
		Run(handler(repo, nil, idService),
			func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
				assert.Equal(t, http.StatusCreated, r.Code)
				assert.Equal(t, 1, len(repo.(*MockDatabase).CreatedResources))
				resource := repo.(*MockDatabase).CreatedResources[0]
				assert.Equal(t, 1, resource.Version())
				assert.Equal(t, "zt570tx3016", resource.ExternalIdentifier())
				assert.NotNil(t, (*resource.Administrative())["created"])
				// assert.Equal(t, "bib12345678", repo.(*MockDatabase).CreatedResources[0].(map[string]interface{})["sourceid"])
			})
}

func TestCreateResourceNoApiKey(t *testing.T) {
	r := gofight.New()
	r.POST("/v1/resource").
		SetJSON(gofight.D{
			"tacoIdentifier": "oo000oo0001",
			"sourceId":       "bib12345678",
			"title":          "My work",
		}).
		Run(handler(nil, nil, nil),
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
		SetJSON(postData("request.json")).
		Run(handler(nil, nil, nil),
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
			"tacoIdentifier": "oo000oo0001",
			"@type":          "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"title":          "My work",
		}).
		Run(handler(nil, nil, nil),
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
			"tacoIdentifier": "oo000oo0001",
			"@context":       "http://example.com", // This is not a valid context
			"@type":          "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":         "world",
			"label":          "My work",
			"preserve":       true,
			"publish":        true,
			"sourceId":       "bib12345678"}).
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
				SetHeader(gofight.H{
					"On-Behalf-Of": "lmcrae@stanford.edu",
				}).
				SetJSON(postData("request.json")).
				Run(handler(NewMockErrorDatabase(nil), nil, nil),
					func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {})
		})
}
