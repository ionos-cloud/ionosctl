package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

type Response struct {
	sdkgo.APIResponse
}

// RegistriesService is a wrapper around ionoscloud.Registry
type RegistriesService interface {
	List(filterName string) (sdkgo.RegistriesResponse, *Response, error)
}

type registriesService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ RegistriesService = &registriesService{}

func NewRegistriesService(client *config.Client, ctx context.Context) RegistriesService {
	return &registriesService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *registriesService) List(filterName string) (sdkgo.RegistriesResponse, *Response, error) {
	req := svc.client.RegistriesApi.RegistriesGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	registryList, res, err := svc.client.RegistriesApi.RegistriesGetExecute(req)
	return registryList, &Response{*res}, err
}
