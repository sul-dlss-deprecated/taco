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
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sul-dlss-labs/taco/config"
	baloo "gopkg.in/h2non/baloo.v3"
	"gopkg.in/h2non/gentleman.v2/plugins/multipart"
)

var apiToken string

func setupTest() *baloo.Client {
	conf := config.NewConfig()
	mySigningKey := []byte(conf.SecretKey)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Date(2099, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
		Issuer:    "test",
		Subject:   "lmcrae@stanford.edu",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(mySigningKey)
	remoteHost, ok := os.LookupEnv("TEST_REMOTE_ENDPOINT")
	if !ok {
		port := config.NewConfig().Port
		remoteHost = fmt.Sprintf("localhost:%v", port)
	}
	return baloo.New(fmt.Sprintf("http://%s", remoteHost)).
		SetHeader("Authorization", fmt.Sprintf("bearer %s", ss))
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

	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	setupTest().Post("/v1/resource").
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

	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err = json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	setupTest().Post("/v1/resource").
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
	patchData["externalIdentifier"] = externalIdentifier
	patchData["tacoIdentifier"] = tacoIdentifier

	setupTest().Patch(fmt.Sprintf("/v1/resource/%s", id)).
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

	byt, err := ioutil.ReadFile("../examples/request.json")
	if err != nil {
		panic(err)
	}
	var postData map[string]interface{}

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}

	setupTest().Post("/v1/resource").
		SetHeader("On-Behalf-Of", "lmcrae@stanford.edu").
		JSON(postData).
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

	file := multipart.FormFile{Name: "upload", Reader: strings.NewReader("data")}
	files := []multipart.FormFile{file}
	setupTest().Post("/v1/file").
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
