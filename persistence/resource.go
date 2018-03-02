package persistence

import (
	"github.com/go-openapi/strfmt"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	ID        string     `json:"id"`
	AtType    strfmt.URI `json:"@type"`
	AtContext strfmt.URI `json:"@context"`
	Access    string     `json:"access"`
	Label     string     `json:"label"`
	Preserve  bool       `json:"sdrPreserve"`
	Publish   bool       `json:"publish"`
	SourceID  string     `json:"source_id"`
}

// NewResource casts parameters into a persisable resource
func NewResource(id string, params *models.Resource) *Resource {
	return &Resource{
		ID:        id,
		Access:    *params.Access,
		AtContext: *params.AtContext,
		AtType:    *params.AtType,
		Label:     *params.Label,
		Preserve:  *params.Preserve,
		Publish:   *params.Publish,
		SourceID:  params.SourceID,
	}
}
