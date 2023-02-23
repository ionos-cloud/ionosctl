package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

// NameService is a contract for the name service.
type NameService interface {
	Get(name string) (*Response, error)
}

// NameServiceOp is an implementation of the NameService interface.
type nameService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ NameService = &nameService{}

// NewNameService returns a new NameService.
func NewNameService(client *config.Client, ctx context.Context) NameService {
	return &nameService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

// Get returns a response.
func (svc *nameService) Get(name string) (*Response, error) {
	req := svc.client.NamesApi.NamesCheckUsage(svc.context, name)
	res, err := svc.client.NamesApi.NamesCheckUsageExecute(req)
	return &Response{*res}, err
}
