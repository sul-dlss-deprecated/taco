package identifier

import (
	"github.com/google/uuid"
	"github.com/sul-dlss-labs/taco/datautils"
)

type uuidService struct{}

// NewUUIDService creates a new instance of the UUID identifier service
func NewUUIDService() Service {
	return &uuidService{}
}

func (d *uuidService) Mint(resource *datautils.Resource) (string, error) {
	resourceID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return resourceID.String(), nil
}
