package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

type RepositoryService interface {
	Delete(regId string, name string) (*Response, error)
}

type repositoryService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ RepositoryService = &repositoryService{}

func NewRepositoryService(client *config.Client, ctx context.Context) RepositoryService {
	return &repositoryService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *repositoryService) Delete(regId string, name string) (*Response, error) {
	req := svc.client.RepositoriesApi.RegistriesRepositoriesDelete(svc.context, regId, name)
	res, err := svc.client.RepositoriesApi.RegistriesRepositoriesDeleteExecute(req)
	return &Response{*res}, err
}
