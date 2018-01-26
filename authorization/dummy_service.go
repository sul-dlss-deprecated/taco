package authorization

import "github.com/sul-dlss-labs/taco/datautils"

type dummyAuthorizationService struct{}

// NewDummyAuthorizationService returns a new instance of the dummy service
func NewDummyAuthorizationService() Service {
	return &dummyAuthorizationService{}
}

func (d *dummyAuthorizationService) CanCreateResourceOfType(agentID string, resourceType string) bool {
	return agentID == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanRetrieveResource(agentID string, res *datautils.Resource) bool {
	return agentID == "lmcrae@stanford.edu"
}

func (d *dummyAuthorizationService) CanDeleteResource(agentID string, id string) bool {
	return agentID == "lmcrae@stanford.edu"
}
