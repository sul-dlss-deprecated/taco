package test

import (
	"flag"
	"fmt"
	"testing"

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

const schema = `{
  "title": "Create response",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    }
  },
  "required": ["id"]
}`

func TestBalooSimple(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	setupTest().Post("/v1/resource").
		SetHeader("Foo", "Bar").
		JSON(map[string]string{
			"@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
			"@type":    "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
			"access":   "world",
			"id":       "oo000oo0001",
			"label":    "My SDR3 resource",
			"preserve": "true",
			"publish":  "true",
			"sourceId": "bib12345678"}).
		Expect(t).
		Status(200).
		Type("json").
		JSONSchema(schema).
		Done()
}
