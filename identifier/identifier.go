package identifier

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
)

// Service the interface for a service that mints identifiers
type Service interface {
	Mint() (string, error)
}

func NewService(config *config.Config) Service {
	if config.IdentifierServiceHost == "" {
		log.Println("IDENTIFIER_SERVICE_HOST is not set, so using UUID service")
		return NewUUIDService()
	} else {
		return NewRemoteIdentifierService(config)
	}
}
