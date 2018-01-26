package authorization

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	remoteService "github.com/sul-dlss-labs/permissions-service/generated/client"
	"github.com/sul-dlss-labs/permissions-service/generated/client/operations"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/datautils"
)

type remoteAuthorizationService struct {
	TransportConfig *remoteService.TransportConfig
}

// NewRemoteAuthorizationService creates a new instance of the authorization service
func NewRemoteAuthorizationService(config *config.Config) Service {
	host := config.AuthorizationServiceHost
	return &remoteAuthorizationService{
		TransportConfig: remoteService.DefaultTransportConfig().WithHost(host),
	}
}

func (d *remoteAuthorizationService) CanCreateResourceOfType(agentID string, resourceType string) bool {
	c := remoteService.NewHTTPClientWithConfig(nil, d.TransportConfig)
	params := paramsFor("create", resourceType)
	ok, err := c.Operations.QueryAction(params, authInfo(agentID))
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) CanRetrieveResource(agentID string, res *datautils.Resource) bool {
	c := remoteService.NewHTTPClientWithConfig(nil, d.TransportConfig)

	params := paramsFor("retrieve", res.ID())
	ok, err := c.Operations.QueryAction(params, authInfo(agentID))
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) CanDeleteResource(agentID string, id string) bool {
	c := remoteService.NewHTTPClientWithConfig(nil, d.TransportConfig)

	params := paramsFor("delete", id)
	ok, err := c.Operations.QueryAction(params, authInfo(agentID))
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func authInfo(agentID string) runtime.ClientAuthInfoWriter {
	return client.APIKeyAuth("On-Behalf-Of", "header", agentID)
}

func paramsFor(action string, resource string) *operations.QueryActionParams {
	return operations.NewQueryActionParams().WithAction(action).WithResource(resource)
}
