package persistence

import (
	"github.com/go-openapi/strfmt"
)

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	ID        string     `json:"id"`
	AtType    strfmt.URI `json:"@type"`
	AtContext strfmt.URI `json:"@context"`
	Access    string     `json:"@access"`
	Label     string     `json:"label"`
	Preserve  bool       `json:"@preserve"`
	Publish   bool       `json:"@publish"`
	SourceID  string     `json:"source_id"`
}
