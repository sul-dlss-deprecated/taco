package identifier

import (
	"github.com/google/uuid"
)

type uuidService struct{}

// NewUUIDService creates a new instance of the UUID identifier service
func NewUUIDService() Service {
	return &uuidService{}
}

func (d *uuidService) Mint() (string, error) {
	resourceID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return resourceID.String(), nil
}
