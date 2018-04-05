package authorization

import (
	"log"

	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

// Service can answer queries about whether an agent can take a specific action
type Service interface {
	CanCreateResourceOfType(agent *Agent, resourceType string) bool
	CanRetrieveResource(agent *Agent, resource *datautils.Resource) bool
	CanDeleteResource(agent *Agent, id string) bool
	CanUpdateResource(agent *Agent, resource *datautils.Resource) bool
}

// NewService creates a new instance of the authorization service
func NewService(config *config.Config) Service {
	if config.AuthorizationServiceHost == "" {
		log.Println("AUTHORIZATION_SERVICE_HOST is not set, so using dummy authorization service")
		return NewDummyAuthorizationService()
	}
	return NewRemoteAuthorizationService(config)
}
