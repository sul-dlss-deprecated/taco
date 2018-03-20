package datautils

import "fmt"

// JSONObject represents a JSON object.
type JSONObject map[string]interface{}

// GetS returns the string value at key
func (d *JSONObject) GetS(key string) string {
	d.ensureKeyExists(key)
	return (*d)[key].(string)
}

// GetI returns the int value at key
func (d *JSONObject) GetI(key string) int {
	d.ensureKeyExists(key)
	return (*d)[key].(int)
}

// GetA returns the array value at key
func (d *JSONObject) GetA(key string) *JSONArray {
	d.ensureKeyExists(key)
	arr := JSONArray((*d)[key].([]interface{}))
	return &arr
}

// GetObj returns the JSONObject value at the given key
func (d *JSONObject) GetObj(key string) *JSONObject {
	d.ensureKeyExists(key)
	obj := JSONObject((*d)[key].(map[string]interface{}))
	return &obj
}

// HasKey returns a boolean value:
// false unless the key exists
func (d *JSONObject) HasKey(key string) bool {
	return (*d)[key] != nil
}

func (d *JSONObject) ensureKeyExists(key string) {
	if d.HasKey(key) {
		return
	}
	panic(fmt.Errorf("No key found for %s", key))
}
