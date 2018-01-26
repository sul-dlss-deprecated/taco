package authorization

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Service can answer queries about whether an agent can take a specific action
type Service interface {
	CanCreateResourceOfType(agentID string, resourceType string) bool
	CanRetrieveResource(agentID string, resource *datautils.Resource) bool
	CanDeleteResource(agentID string, id string) bool
}

// NewService creates a new instance of the authorization service
func NewService(config *config.Config) Service {
	if config.AuthorizationServiceHost == "" {
		log.Println("AUTHORIZATION_SERVICE_HOST is not set, so using dummy authorization service")
		return NewDummyAuthorizationService()
	}
	return NewRemoteAuthorizationService(config)
}
