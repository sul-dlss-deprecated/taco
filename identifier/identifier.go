package identifier

import "github.com/google/uuid"

// Service the interface for a service that mints identifiers
type Service interface {
	Mint() (string, error)
}
type uuidIdentifierService struct{}

// NewService creates a new instance of the identifier service
// TODO: This should create a DRUID service
func NewService() Service {
	return &uuidIdentifierService{}
}

func (d *uuidIdentifierService) Mint() (string, error) {
	resourceID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return resourceID.String(), nil
}
