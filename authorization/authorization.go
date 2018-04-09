package authorization

import (
	"github.com/sul-dlss-labs/taco/datautils"
)

// Service can answer queries about whether an agent can take a specific action
type Service interface {
	CanCreateResourceOfType(string) bool
	CanRetrieveResource(*datautils.Resource) bool
	CanUpdateResource(*datautils.Resource) bool
}

type dummyAuthorizationService struct {
	agent *Agent
}

// NewService creates a new instance of the authorization service
func NewService(agent *Agent) Service {
	return &dummyAuthorizationService{agent: agent}
}

func (d *dummyAuthorizationService) CanCreateResourceOfType(resourceType string) bool {
	return d.agent.Identifier == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanRetrieveResource(res *datautils.Resource) bool {
	return d.agent.Identifier == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanUpdateResource(res *datautils.Resource) bool {
	return d.agent.Identifier == "lmcrae@stanford.edu"
}
