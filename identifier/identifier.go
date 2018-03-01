package identifier

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
)

// Service the interface for a service that mints identifiers
type Service interface {
	Mint() (string, error)
}

// NewService builds the identifier service
func NewService(config *config.Config) *TypeAwareService {
	externalService := externalServiceOrFallback(config.IdentifierServiceHost)
	return &TypeAwareService{
		UUIDService:       NewUUIDService(),
		IdentifierService: externalService,
	}
}

func externalServiceOrFallback(host string) Service {
	if host == "" {
		log.Println("IDENTIFIER_SERVICE_HOST is not set, so using UUID service")
		return NewUUIDService()
	}
	return NewRemoteIdentifierService(host)
}
