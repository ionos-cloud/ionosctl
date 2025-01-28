package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
)

// NameService is a contract for the name service.
type NameService interface {
	Head(name string) (*shared.APIResponse, error)
}

// NameServiceOp is an implementation of the NameService interface.
type nameService struct {
	client  *containerregistry.APIClient
	context context.Context
}

var _ NameService = &nameService{}

// NewNameService returns a new NameService.
func NewNameService(client *client2.Client, ctx context.Context) NameService {
	return &nameService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

// Get returns a *shared.APIResponse.
func (svc *nameService) Head(name string) (*shared.APIResponse, error) {
	req := svc.client.NamesApi.NamesCheckUsage(svc.context, name)
	res, err := svc.client.NamesApi.NamesCheckUsageExecute(req)
	return res, err
}
