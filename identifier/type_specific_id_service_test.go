package identifier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sul-dlss-labs/taco/datautils"
)

type FakeService struct {
	returns string
}

func (d *FakeService) Mint(*datautils.Resource) (string, error) { return d.returns, nil }

func TestTypeSpecificIdServiceCollection(t *testing.T) {
	localService := &FakeService{returns: "local"}
	remoteService := &FakeService{returns: "remote"}
	service := &TypeSpecificIDService{
		localService:  localService,
		remoteService: remoteService,
	}
	collectionJSON := datautils.JSONObject{"@type": datautils.CollectionTypes[0]}
	id, _ := service.Mint(datautils.NewResource(collectionJSON))
	assert.Equal(t, "remote", id)

	objectJSON := datautils.JSONObject{"@type": datautils.ObjectTypes[0]}
	id, _ = service.Mint(datautils.NewResource(objectJSON))
	assert.Equal(t, "remote", id)

	filesetJSON := datautils.JSONObject{"@type": datautils.FilesetType}
	id, _ = service.Mint(datautils.NewResource(filesetJSON))
	assert.Equal(t, "local", id)
}
