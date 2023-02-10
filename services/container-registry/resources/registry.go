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
	Get(id string) (sdkgo.RegistryResponse, *Response, error)
	List(filterName string) (sdkgo.RegistriesResponse, *Response, error)
	Post(input sdkgo.PostRegistryInput) (sdkgo.PostRegistryOutput, *Response, error)
	Delete(id string) (*Response, error)
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

func (svc *registriesService) Post(input sdkgo.PostRegistryInput) (sdkgo.PostRegistryOutput, *Response, error) {
	req := svc.client.RegistriesApi.RegistriesPost(svc.context).PostRegistryInput(input)
	registryList, res, err := svc.client.RegistriesApi.RegistriesPostExecute(req)
	return registryList, &Response{*res}, err
}

func (svc *registriesService) Get(id string) (sdkgo.RegistryResponse, *Response, error) {
	req := svc.client.RegistriesApi.RegistriesFindById(svc.context, id)
	reg, res, err := svc.client.RegistriesApi.RegistriesFindByIdExecute(req)
	return reg, &Response{*res}, err
}

func (svc *registriesService) Delete(id string) (*Response, error) {
	req := svc.client.RegistriesApi.RegistriesDelete(svc.context, id)
	res, err := svc.client.RegistriesApi.RegistriesDeleteExecute(req)
	return &Response{*res}, err
}
