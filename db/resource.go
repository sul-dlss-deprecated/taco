package db

import (
	"fmt"
)

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource map[string]interface{}

// NewResource creates a new resource instance
func NewResource(data map[string]interface{}) Resource {
	return Resource(data)
}

// ID returns the documents identifier
func (d *Resource) ID() string {
	return d.GetS("id")
}

// GetS returns the string value at key
func (d *Resource) GetS(key string) string {
	if (*d)[key] == nil {
		panic(fmt.Errorf("No key found for %s", key))
	}
	val := (*d)[key].(string)
	return val
}

func (d *Resource) String() string {
	return fmt.Sprintf("<Resource id: '%s'>", d.ID())
}
