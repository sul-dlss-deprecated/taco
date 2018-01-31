package test

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"testing"

	"github.com/buger/jsonparser"
	baloo "gopkg.in/h2non/baloo.v3"
)

var port string

// For a custom port, invoke this as:
//   go test test/integration_test.go --port 3000
func init() {
	flag.StringVar(&port, "port", "8080", "port for test server")
}

func setupTest() *baloo.Client {
	return baloo.New(fmt.Sprintf("http://localhost:%s", port))
}

const createSchema = `{
  "title": "Create response",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    }
  },
  "required": ["id"]
}`

var id string

func TestBalooSimple(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	postData := map[string]interface{}{
		"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
		"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
		"access":   "world",
		"id":       "oo000oo0001",
		"label":    "My SDR3 resource",
		"preserve": true,
		"publish":  true,
		"sourceId": "bib12345678"}

	setupTest().Post("/v1/resource").
		SetHeader("Foo", "Bar").
		JSON(postData).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(createSchema).
		AssertFunc(assert).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

// assert implements an assertion function with custom validation logic.
// If the assertion fails it should return an error.
// This has the side effect of setting the top level id variable
// which we use for making a subsequent request.
func assert(res *http.Response, req *http.Request) error {
	// TODO: this parsing would be unnecessary if we had a Location header
	//       Then we could just do res.Location()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	jsonID, _ := jsonparser.GetString(buf.Bytes(), "id")
	id = jsonID
	return nil
}
