package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	containerregistry "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
)

type RepositoryService interface {
	Delete(regId string, name string) (*containerregistry.APIResponse, error)
}

type repositoryService struct {
	client  *containerregistry.APIClient
	context context.Context
}

var _ RepositoryService = &repositoryService{}

func NewRepositoryService(client *client2.Client, ctx context.Context) RepositoryService {
	return &repositoryService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *repositoryService) Delete(regId string, name string) (*containerregistry.APIResponse, error) {
	req := svc.client.RepositoriesApi.RegistriesRepositoriesDelete(svc.context, regId, name)
	res, err := svc.client.RepositoriesApi.RegistriesRepositoriesDeleteExecute(req)
	return res, err
}
