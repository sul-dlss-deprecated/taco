package validators

import (
	"github.com/sul-dlss-labs/taco/datautils"
	"github.com/sul-dlss-labs/taco/generated/models"
)

// ResourceValidator is the interface for validators that check the resources format
type ResourceValidator interface {
	ValidateResource(resource *datautils.Resource) *models.ErrorResponseErrors
}
