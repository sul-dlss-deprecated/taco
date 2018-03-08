package persistence

import (
	"github.com/sul-dlss-labs/taco/generated/models"
)

// Resource represents the resource as it exists in the persistence layer
// this is very similar to models.Resource, but COULD vary, so we should
// keep them separated
type Resource struct {
	ID             string                        `json:"id"`
	AtType         string                        `json:"@type"`
	AtContext      string                        `json:"@context"`
	Access         models.ResourceAccess         `json:"access"`
	Administrative models.ResourceAdministrative `json:"administrative"`
	Identification models.ResourceIdentification `json:"identification"`
	Label          string                        `json:"label"`
	//Preserve       bool                          `json:"sdrPreserve"`
	//Publish        bool                          `json:"publish"`
	//SourceID string `json:"source_id"`
}
