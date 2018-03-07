package persistence

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource gin.H

// NewResource creates a new resource instance
func NewResource(data gin.H) Resource {
	return Resource(data)
}

// ID returns the documents identifier
func (d *Resource) ID() string {
	return d.GetS("id")
}

// GetS returns the string value at key
func (d *Resource) GetS(key string) string {
	return (*d)[key].(string)
}

func (d *Resource) String() string {
	return fmt.Sprintf("<Resource id: '%s'>", d.ID())
}
