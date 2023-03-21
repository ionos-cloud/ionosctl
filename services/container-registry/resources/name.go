package resources

import (
	"context"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

// NameService is a contract for the name service.
type NameService interface {
	Head(name string) (*sdkgo.APIResponse, error)
}

// NameServiceOp is an implementation of the NameService interface.
type nameService struct {
	client  *sdkgo.APIClient
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

// Get returns a sdkgo.APIResponse.
func (svc *nameService) Head(name string) (*sdkgo.APIResponse, error) {
	req := svc.client.NamesApi.NamesCheckUsage(svc.context, name)
	res, err := svc.client.NamesApi.NamesCheckUsageExecute(req)
	return res, err
}
