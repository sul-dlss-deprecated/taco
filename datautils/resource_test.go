package datautils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithID(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithID("8888")
	assert.Equal(t, "8888", resource.ID())
}

func TestWithExternalIdentifier(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithExternalIdentifier("8888")
	assert.Equal(t, "8888", resource.JSON.GetS("externalIdentifier"))
}

func TestWithVersion(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithVersion(5)
	assert.Equal(t, 5, resource.Version())
}

func TestWithLabel(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithLabel("My file")
	assert.Equal(t, "My file", resource.Label())
}

func TestWithFileLocation(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithFileLocation("s3://bucket/key")
	assert.Equal(t, "s3://bucket/key", resource.FileLocation())
}

func TestWithMimeType(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithMimeType("text/plain")
	assert.Equal(t, "text/plain", resource.MimeType())
}

func TestWithType(t *testing.T) {
	resource := NewResource(JSONObject{})
	resource.WithType("Foo")
	assert.Equal(t, "Foo", resource.Type())
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
	json := JSONObject{}
	resource := NewResource(json).WithType("Bar")
	assert.Equal(t, "<Resource @type:'Bar'>", resource.String())
	resource = resource.WithID("Foo")
	assert.Equal(t, "<Resource id:'Foo' @type:'Bar'>", resource.String())
}
