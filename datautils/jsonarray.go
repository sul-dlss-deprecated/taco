package datautils

// JSONArray represents a JSON array.
type JSONArray []interface{}

// GetS returns the equivalent string array
func (d *JSONArray) GetS() []string {
	b := make([]string, len(*d), len(*d))
	for i := range *d {
		b[i] = (*d)[i].(string)
	}
	return b
}

// GetObj returns the equivalent array of JSONObjects
func (d *JSONArray) GetObj() []*JSONObject {
	b := make([]*JSONObject, len(*d), len(*d))
	for i := range *d {
		obj := JSONObject((*d)[i].(map[string]interface{}))
		b[i] = &obj

	}
	return b
}
