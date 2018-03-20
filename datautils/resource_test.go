package datautils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	json := JSONObject{"@type": "Foo"}
	resource := NewResource(json)
	assert.Equal(t, "Foo", resource.Type())
}

func TestID(t *testing.T) {
	json := JSONObject{"id": "123"}
	resource := NewResource(json)
	assert.Equal(t, "123", resource.ID())
}

func TestIsObject(t *testing.T) {
	fileJSON := JSONObject{"@type": FileType}
	file := NewResource(fileJSON)
	assert.False(t, file.IsObject())

	objJSON := JSONObject{"@type": "http://sdr.sul.stanford.edu/models/sdr3-document.jsonld"}
	obj := NewResource(objJSON)
	assert.True(t, obj.IsObject())

	mapJSON := JSONObject{"@type": "http://sdr.sul.stanford.edu/models/sdr3-map.jsonld"}
	objmap := NewResource(mapJSON)
	assert.True(t, objmap.IsObject())
}

func TestIsCollection(t *testing.T) {
	fileJSON := JSONObject{"@type": FileType}
	file := NewResource(fileJSON)
	assert.False(t, file.IsCollection())

	colJSON := JSONObject{"@type": "http://sdr.sul.stanford.edu/models/sdr3-collection.jsonld"}
	col := NewResource(colJSON)
	assert.True(t, col.IsCollection())

	subcolJSON := JSONObject{"@type": "http://sdr.sul.stanford.edu/models/sdr3-exhibit.jsonld"}
	subcol := NewResource(subcolJSON)
	assert.True(t, subcol.IsCollection())
}

func TestIsFile(t *testing.T) {
	notFileJSON := JSONObject{"@type": "Foo"}
	notFile := NewResource(notFileJSON)
	assert.False(t, notFile.IsFile())

	fileJSON := JSONObject{"@type": FileType}
	file := NewResource(fileJSON)
	assert.True(t, file.IsFile())
}

func TestString(t *testing.T) {
	json := JSONObject{"id": "Foo"}
	resource := NewResource(json)
	assert.Equal(t, "<Resource id: 'Foo'>", resource.String())
}
