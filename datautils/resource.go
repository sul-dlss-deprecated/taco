package datautils

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

// ID returns the document's identifier
func (d *Resource) ID() string {
	return d.GetS("id")
}

// Version returns the document's version
func (d *Resource) Version() int {
	return d.GetI("version")
}

// GetS returns the string value at key
func (d *Resource) GetS(key string) string {
	if (*d)[key] == nil {
		panic(fmt.Errorf("No key found for %s", key))
	}
	return (*d)[key].(string)
}

// GetI returns the int value at key
func (d *Resource) GetI(key string) int {
	if (*d)[key] == nil {
		panic(fmt.Errorf("No key found for %s", key))
	}
	return (*d)[key].(int)
}

func (d *Resource) String() string {
	return fmt.Sprintf("<Resource id: '%s'>", d.ID())
}
