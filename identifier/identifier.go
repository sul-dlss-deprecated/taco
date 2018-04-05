package identifier

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Service the interface for a service that mints identifiers
type Service interface {
	Mint(resource *datautils.Resource) (string, error)
}

func NewService(config *config.Config) Service {
	return &TypeSpecificIDService{
		remoteService: remoteService(config),
		localService:  NewUUIDService(),
	}
}

func remoteService(config *config.Config) Service {
	if config.IdentifierServiceHost == "" {
		log.Println("IDENTIFIER_SERVICE_HOST is not set, so using UUID service")
		return NewUUIDService()
	} else {
		return NewRemoteIdentifierService(config)
	}
}
