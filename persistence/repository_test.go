package persistence

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/sessionbuilder"
)

func TestSaveAndRetrieve(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	repo := initRepo()

	resource := map[string]*dynamodb.AttributeValue{}
	resource[PrimaryKey] = &dynamodb.AttributeValue{S: &id}
	label := "Hello world"
	resource["Label"] = &dynamodb.AttributeValue{S: &label}

	err := repo.CreateItem(resource)
	assert.Nil(t, err)
	item, err := repo.GetByID(id)
	assert.Nil(t, err)
	assert.Equal(t, *item.Label, label)
}

func initRepo() *DynamoRepository {
	config := config.NewConfig()
	conn := db.NewConnection(config, sessionbuilder.NewAwsSession(config))
	return NewDynamoRepository(config, conn)
}
