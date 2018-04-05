package authorization

import "github.com/sul-dlss-labs/taco/datautils"

type dummyAuthorizationService struct{}

// NewDummyAuthorizationService returns a new instance of the dummy service
func NewDummyAuthorizationService() Service {
	return &dummyAuthorizationService{}
}

func (d *dummyAuthorizationService) CanCreateResourceOfType(agent *Agent, resourceType string) bool {
	return agent.Identifier == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanRetrieveResource(agent *Agent, res *datautils.Resource) bool {
	return agent.Identifier == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanUpdateResource(agent *Agent, res *datautils.Resource) bool {
	return agent.Identifier == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanDeleteResource(agent *Agent, id string) bool {
	return agent.Identifier == "lmcrae@stanford.edu"
}
