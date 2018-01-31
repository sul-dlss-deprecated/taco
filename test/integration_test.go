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
		JSON(map[string]string{"title": "value1", "sourceId": "value2"}).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(schema).
		Done()
}
