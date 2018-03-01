package identifier

import (
	"github.com/go-openapi/strfmt"
	"github.com/sul-dlss-labs/taco/persistence"
)

// TypeAwareService is a identifier service that returns the appropriate type
// of identifier for the object
type TypeAwareService struct {
	UUIDService       Service
	IdentifierService Service
}

// Mint creates the identifiers for the new object based on the resource type.
// If the resource is a DRO or a Collection, it will return a DRUID from the
// identifier-service (if available, otherwise a uuid) and secondly will return
// a sdrIdentifier (also a uuid)
func (d *TypeAwareService) Mint(resourceType strfmt.URI) (string, string, error) {
	sdrUUID, err := d.UUIDService.Mint()
	if err != nil {
		return "", "", err
	}
	if persistence.IsObjectOrCollection(resourceType) {
		druidOrUUID, err := d.IdentifierService.Mint()
		if err != nil {
			return "", "", err
		}
		return sdrUUID, druidOrUUID, nil
	}
	return sdrUUID, sdrUUID, nil
}
