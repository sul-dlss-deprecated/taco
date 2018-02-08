package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/db"
)

func TestSaveAndRetrieve(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	id := "9999"
	config := config.NewConfig()
	repo := NewDynamoRepository(config, db.NewConnection(config))
	resource := &Resource{ID: id, Label: "Hello world"}
	err := repo.CreateItem(resource)
	assert.Nil(t, err)
	item, err := repo.GetByID(id)
	assert.Nil(t, err)
	assert.Equal(t, item.Label, "Hello world")
}
