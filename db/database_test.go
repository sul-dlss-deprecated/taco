package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/aws_session"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

func TestSaveAndRetrieve(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	database := initDatabase()
	resource := datautils.Resource{
		"id":    id,
		"label": "Hello world",
	}
	err := database.Insert(resource)
	assert.Nil(t, err)
}

func initDatabase() *DynamodbDatabase {
	testConfig := config.NewConfig()
	return &DynamodbDatabase{
		Connection: Connect(aws_session.Connect(testConfig.AwsDisableSSL), testConfig.DynamodbEndpoint),
		Table:      testConfig.ResourceTableName,
	}
}
