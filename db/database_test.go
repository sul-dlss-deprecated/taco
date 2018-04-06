package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/aws_session"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestRetrieveVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	database := initDatabase()
	id := "9999"
	json := datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777777",
		"version":    1,
		"label":      "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}

	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777778",
		"version":    2,
		"label":      "Middle one",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777779",
		"version":    3,
		"label":      "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	version := "2"
	record, err := database.RetrieveVersion(id, &version)
	assert.Nil(t, err)
	assert.Equal(t, "Middle one", record.JSON.GetS("label"))
}

func TestRetrieveLatest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	database := initDatabase()
	id := "9999"
	json := datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777777",
		"version":    1,
		"label":      "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}

	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777778",
		"version":    2,
		"label":      "Middle one",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	json = datautils.JSONObject{
		"externalIdentifier": id,
		"tacoIdentifier":         "7777779",
		"version":    3,
		"label":      "Hello world",
	}
	if err := database.Insert(datautils.NewResource(json)); err != nil {
		panic(err)
	}
	record, err := database.RetrieveLatest(id)
	assert.Nil(t, err)
	assert.Equal(t, 3, record.Version())
}

func TestRetrieveLatestNotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	_, err := initDatabase().RetrieveLatest("8888")
	assert.NotNil(t, err)
}

func initDatabase() Database {
	testConfig := config.NewConfig()
	return &DynamodbDatabase{
		Connection: Connect(aws_session.Connect(testConfig.AwsDisableSSL), testConfig.DynamodbEndpoint),
		Table:      testConfig.ResourceTableName,
	}
}
