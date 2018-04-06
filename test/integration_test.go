package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/buger/jsonparser"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
	baloo "gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

func setupTest() *baloo.Client {
	remoteHost, ok := os.LookupEnv("TEST_REMOTE_ENDPOINT")
	if !ok {
		port := config.NewConfig().Port
		remoteHost = fmt.Sprintf("localhost:%v", port)
	}
	return baloo.New(fmt.Sprintf("http://%s", remoteHost))
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

// ID variables to verify updates
var externalIdentifier string
var tacoIdentifier string

func TestCreateAndDestroyResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(fromFile("request.json")).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(200).
		Type("json").
		Done()

	setupTest().Delete(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(204).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(404).
		Done()

}

func TestUpdateResource(t *testing.T) {
	if testing.Short() {
		t.Skip("skpping integration test in short mode")
	}

	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(fromFile("request.json")).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		AddQuery("version", "1").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assertRetrieveResponse).
		Done()

	// This is required to ensure we're not passing an example with
	// incorrect identifiers
	patchData := fromFile("update_request.json")
	(*patchData)["externalIdentifier"] = externalIdentifier
	(*patchData)["tacoIdentifier"] = tacoIdentifier

	setupTest().Patch(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		SetHeader("Content-Type", "application/json").
		JSON(patchData).
		Expect(t).
		Status(200).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	time.Sleep(5 * time.Millisecond)

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assertUpdatedResourceResponse).
		Done()
}

func TestRetrieveVersions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(fromFile("request.json")).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		AddQuery("version", "1").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(200).
		Type("json").
		Done()
}

func TestCreateFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// First create the DRO to attach it to.
	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(fromFile("deposit_object.json")).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	fileset := fromFile("deposit_fileset.json")
	(*fileset.GetObj("structural"))["isContainedBy"] = id
	// Next, create the fileset to attach it to.
	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(fileset).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	file := multipart.FormFile{Name: "upload", Reader: strings.NewReader("sample data")}
	files := []multipart.FormFile{file}
	filePath := fmt.Sprintf("/v1/resource/%s/file", id)
	setupTest().Post(filePath).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Files(files).
		Expect(t).
		Status(201).
		Type("json").
		JSONSchema(resourceSchema).
		AssertFunc(assertResourceResponse).
		Done()

	setupTest().Get(fmt.Sprintf("/v1/resource/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		Status(200).
		Type("json").
		Done()

	setupTest().Get(fmt.Sprintf("/v1/file/%s", id)).
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		Expect(t).
		BodyEquals("sample data").
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

func assertRetrieveResponse(res *http.Response, req *http.Request) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	externalIdentifier, _ = jsonparser.GetString(buf.Bytes(), "externalIdentifier")
	tacoIdentifier, _ = jsonparser.GetString(buf.Bytes(), "tacoIdentifier")
	return nil
}

func assertUpdatedResourceResponse(res *http.Response, req *http.Request) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	jsonLabel, _ := jsonparser.GetString(buf.Bytes(), "label")
	if jsonLabel != "UPDATED: Leon Kolb Collection of Portraits" {
		return errors.New("UpdateResource failure")
	}
	return nil
}

func fromFile(file string) *datautils.JSONObject {
	byt, err := ioutil.ReadFile(fmt.Sprintf("../examples/%s", file))
	if err != nil {
		panic(err)
	}
	var postData datautils.JSONObject

	if err = json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}
	return &postData
}
