package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// RegistriesService is a wrapper around ionoscloud.Registry
type RegistriesService interface {
	Get(id string) (sdkgo.RegistryResponse, *shared.APIResponse, error)
	List(filterName string) (sdkgo.RegistriesResponse, *shared.APIResponse, error)
	Post(input sdkgo.PostRegistryInput) (sdkgo.PostRegistryOutput, *shared.APIResponse, error)
	Delete(id string) (*shared.APIResponse, error)
	Patch(id string, input sdkgo.PatchRegistryInput) (sdkgo.RegistryResponse, *shared.APIResponse, error)
	Put(id string, input sdkgo.PutRegistryInput) (sdkgo.PutRegistryOutput, *shared.APIResponse, error)
}

type registriesService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ RegistriesService = &registriesService{}

func NewRegistriesService(client *client2.Client, ctx context.Context) RegistriesService {
	return &registriesService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *registriesService) List(filterName string) (sdkgo.RegistriesResponse, *shared.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	registryList, res, err := svc.client.RegistriesApi.RegistriesGetExecute(req)
	return registryList, res, err
}

func (svc *registriesService) Post(input sdkgo.PostRegistryInput) (sdkgo.PostRegistryOutput, *shared.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesPost(svc.context).PostRegistryInput(input)
	registryList, res, err := svc.client.RegistriesApi.RegistriesPostExecute(req)
	return registryList, res, err
}

func (svc *registriesService) Get(id string) (sdkgo.RegistryResponse, *shared.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesFindById(svc.context, id)
	reg, res, err := svc.client.RegistriesApi.RegistriesFindByIdExecute(req)
	return reg, res, err
}

func (svc *registriesService) Delete(id string) (*shared.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesDelete(svc.context, id)
	res, err := svc.client.RegistriesApi.RegistriesDeleteExecute(req)
	return res, err
}

func (svc *registriesService) Patch(id string, input sdkgo.PatchRegistryInput) (
	sdkgo.RegistryResponse, *shared.APIResponse, error,
) {
	req := svc.client.RegistriesApi.RegistriesPatch(svc.context, id).PatchRegistryInput(input)
	reg, res, err := svc.client.RegistriesApi.RegistriesPatchExecute(req)
	return reg, res, err
}

func (svc *registriesService) Put(id string, input sdkgo.PutRegistryInput) (
	sdkgo.PutRegistryOutput, *shared.APIResponse, error,
) {
	req := svc.client.RegistriesApi.RegistriesPut(svc.context, id).PutRegistryInput(input)
	reg, res, err := svc.client.RegistriesApi.RegistriesPutExecute(req)
	return reg, res, err
}
