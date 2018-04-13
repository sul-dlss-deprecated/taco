package db

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/aws_session"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

var database = initDatabase()

func TestRetrieveVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	resource := datautils.NewResource(datautils.JSONObject{}).
		WithVersion(1).
		WithLabel("Hello world").
		WithExternalIdentifier(id).
		WithID("7777777")

	if err := database.Insert(resource); err != nil {
		panic(err)
	}

	resource = datautils.NewResource(datautils.JSONObject{}).
		WithVersion(2).
		WithLabel("Middle one").
		WithExternalIdentifier(id).
		WithID("7777778")

	if err := database.Insert(resource); err != nil {
		panic(err)
	}

	resource = datautils.NewResource(datautils.JSONObject{}).
		WithVersion(3).
		WithLabel("Hello world").
		WithExternalIdentifier(id).
		WithID("7777779")

	if err := database.Insert(resource); err != nil {
		panic(err)
	}
	version := "2"
	record, err := database.RetrieveVersion(id, &version)
	assert.Nil(t, err)
	assert.Equal(t, "Middle one", record.Label())
}

func TestRetrieveLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	resource := datautils.NewResource(datautils.JSONObject{}).
		WithVersion(1).
		WithLabel("Hello world").
		WithExternalIdentifier(id).
		WithID("7777777")
	if err := database.Insert(resource); err != nil {
		panic(err)
	}

	resource = datautils.NewResource(datautils.JSONObject{}).
		WithVersion(2).
		WithLabel("Middle one").
		WithExternalIdentifier(id).
		WithID("7777778")

	if err := database.Insert(resource); err != nil {
		panic(err)
	}

	resource = datautils.NewResource(datautils.JSONObject{}).
		WithVersion(3).
		WithLabel("Hello world").
		WithExternalIdentifier(id).
		WithID("7777779")

	if err := database.Insert(resource); err != nil {
		panic(err)
	}
	result, err := database.RetrieveLatest(id)
	assert.Nil(t, err)
	assert.Equal(t, 3, result.Version())
}

func TestRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "7777777"
	jsonData := jsonData()
	jsonData["externalIdentifier"] = id
	resource := datautils.NewResource(jsonData)
	err := database.Insert(resource)
	assert.Nil(t, err)

	result, err := database.RetrieveLatest(id)
	assert.Nil(t, err)
	// Data that comes out should be the same as the data that went in.
	assert.Equal(t, jsonData, result.JSON)
}

func TestRetrieveLatestNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_, err := database.RetrieveLatest("8888")
	assert.Equal(t, err.Error(), "Unable to find record for 8888")
}

func TestRetrieveVersionNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	version := "9"
	_, err := database.RetrieveVersion("8888", &version)
	assert.Equal(t, err.Error(), "Unable to find record for 8888 with version: 9")
}

func TestDeleteVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	json := datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":     "7777777",
		"version":            1,
		"label":              "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}

	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":     "7777778",
		"version":            2,
		"label":              "Middle one",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":     "7777779",
		"version":            3,
		"label":              "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	err := database.DeleteVersion("7777779")
	assert.Nil(t, err)

	record, err := database.RetrieveLatest(id)
	assert.Nil(t, err)

	// Version 3 was deleted, so version 2 is now the most recent
	assert.Equal(t, 2, record.Version())
}

func initDatabase() Database {
	testConfig := config.NewConfig()
	return &DynamodbDatabase{
		Connection: Connect(aws_session.Connect(testConfig.AwsDisableSSL), testConfig.DynamodbEndpoint),
		Table:      testConfig.ResourceTableName,
	}
}

func jsonData() datautils.JSONObject {
	byt, err := ioutil.ReadFile("../examples/update_request.json")
	if err != nil {
		panic(err)
	}
	var postData datautils.JSONObject

	if err := json.Unmarshal(byt, &postData); err != nil {
		panic(err)
	}
	return postData
}
