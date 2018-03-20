package datautils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetS(t *testing.T) {
	json := JSONObject{"@type": "Foo"}
	assert.Equal(t, "Foo", json.GetS("@type"))
}

func TestGetObj(t *testing.T) {
	json := JSONObject{"structural": map[string]interface{}{"isContainedBy": "druid:bq582kh2487_fs1"}}
	containedBy := json.GetObj("structural").GetS("isContainedBy")
	assert.Equal(t, "druid:bq582kh2487_fs1", containedBy)
}

func TestGetA(t *testing.T) {
	json := JSONObject{"members": []interface{}{"one", "two"}}
	result := json.GetA("members").GetS()
	assert.Equal(t, "one", result[0])
	assert.Equal(t, "two", result[1])
}

func TestHasKey(t *testing.T) {
	json := JSONObject{"foo": "123"}
	assert.True(t, json.HasKey("foo"))
	assert.False(t, json.HasKey("structural"))
}
