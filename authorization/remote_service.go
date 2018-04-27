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

func (d *remoteAuthorizationService) CanCreateResourceOfType(agent *Agent, resourceType string) bool {
	params := paramsFor("create", resourceType)
	ok, err := d.query(agent, params)
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) CanUpdateResource(agent *Agent, res *datautils.Resource) bool {
	params := paramsFor("update", res.ID())
	ok, err := d.query(agent, params)
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) CanDeleteResource(agent *Agent, id string) bool {
	params := paramsFor("update", id)
	ok, err := d.query(agent, params)
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) CanRetrieveResource(agent *Agent, res *datautils.Resource) bool {
	params := paramsFor("retrieve", res.ID())
	ok, err := d.query(agent, params)
	if err != nil {
		panic(err)
	}
	return ok.Payload.Authorized
}

func (d *remoteAuthorizationService) query(agent *Agent, params *operations.QueryActionParams) (*operations.QueryActionOK, error) {
	c := remoteService.NewHTTPClientWithConfig(nil, d.TransportConfig)

	return c.Operations.QueryAction(params, authInfo(agent))
}

func authInfo(agent *Agent) runtime.ClientAuthInfoWriter {
	return client.APIKeyAuth("On-Behalf-Of", "header", agent.Identifier)
}

func paramsFor(action string, resource string) *operations.QueryActionParams {
	return operations.NewQueryActionParams().WithAction(action).WithResource(resource)
}
