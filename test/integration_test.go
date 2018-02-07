package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/buger/jsonparser"
	"github.com/sul-dlss-labs/taco/config"
	baloo "gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

func setupTest() *baloo.Client {
	port := config.NewConfig().Port
	return baloo.New(fmt.Sprintf("http://localhost:%v", port))
}

const resourceSchema = `{
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
var label string

func TestCreateResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	setupTest().Post("/v1/resource").
		SetHeader("Foo", "Bar").
		JSON(postData).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func TestUpdateResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skpping integration test in short mode")
	}

	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	setupTest().Post("/v1/resource").
		SetHeader("Foo", "Bar").
		JSON(postData).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	byt, err = ioutil.ReadFile("../examples/update_request.json")
	if err != nil {
		panic(err)
	}

	var patchData map[string]interface{}

	if err := json.Unmarshal(byt, &patchData); err != nil {
		panic(err)
	}

	setupTest().Patch(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("Content-Type", "application/json").
		JSON(patchData).
		Expect(t).
		Status(200).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assertUpdatedResourceResponse).
		Done()

		// minimal integration test for label (the only field changed in our update statement)
	if label != patchData["label"] {
		t.Errorf("Updated label not applied.")
	}
}

func TestCreateFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	file := multipart.FormFile{Name: "upload", Reader: strings.NewReader("data")}
	files := []multipart.FormFile{file}
	setupTest().Post("/v1/file").
		SetHeader("Foo", "Bar").
		Files(files).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func TestHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	setupTest().Get("/v1/healthcheck").
		Expect(t).
		Status(200).
		Type("json").
		JSON(map[string]string{"status": "OK"}).
		Done()
}

// assert implements an assertion function with custom validation logic.
// If the assertion fails it should return an error.
// This has the side effect of setting the top level id variable
// which we use for making a subsequent request.
func assertResourceResponse(res *http.Response, req *http.Request) error {
	// TODO: this parsing would be unnecessary if we had a Location header
	//       Then we could just do res.Location()
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	jsonID, _ := jsonparser.GetString(buf.Bytes(), "id")
	id = jsonID
	return nil
}

func assertUpdatedResourceResponse(res *http.Response, req *http.Request) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	jsonID, _ := jsonparser.GetString(buf.Bytes(), "id")
	jsonLabel, _ := jsonparser.GetString(buf.Bytes(), "label")
	id = jsonID
	label = jsonLabel
	return nil
}
