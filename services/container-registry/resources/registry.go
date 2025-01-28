package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	containerregistry "github.com/ionos-cloud/sdk-go-container-registry"
)

// RegistriesService is a wrapper around ionoscloud.Registry
type RegistriesService interface {
	Get(id string) (containerregistry.RegistryResponse, *containerregistry.APIResponse, error)
	List(filterName string) (containerregistry.RegistriesResponse, *containerregistry.APIResponse, error)
	Post(input containerregistry.PostRegistryInput) (containerregistry.PostRegistryOutput, *containerregistry.APIResponse, error)
	Delete(id string) (*containerregistry.APIResponse, error)
	Patch(id string, input containerregistry.PatchRegistryInput) (containerregistry.RegistryResponse, *containerregistry.APIResponse, error)
	Put(id string, input containerregistry.PutRegistryInput) (containerregistry.PutRegistryOutput, *containerregistry.APIResponse, error)
}

type registriesService struct {
	client  *containerregistry.APIClient
	context context.Context
}

var _ RegistriesService = &registriesService{}

func NewRegistriesService(client *client2.Client, ctx context.Context) RegistriesService {
	return &registriesService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *registriesService) List(filterName string) (containerregistry.RegistriesResponse, *containerregistry.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	registryList, res, err := svc.client.RegistriesApi.RegistriesGetExecute(req)
	return registryList, res, err
}

func (svc *registriesService) Post(input containerregistry.PostRegistryInput) (containerregistry.PostRegistryOutput, *containerregistry.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesPost(svc.context).PostRegistryInput(input)
	registryList, res, err := svc.client.RegistriesApi.RegistriesPostExecute(req)
	return registryList, res, err
}

func (svc *registriesService) Get(id string) (containerregistry.RegistryResponse, *containerregistry.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesFindById(svc.context, id)
	reg, res, err := svc.client.RegistriesApi.RegistriesFindByIdExecute(req)
	return reg, res, err
}

func (svc *registriesService) Delete(id string) (*containerregistry.APIResponse, error) {
	req := svc.client.RegistriesApi.RegistriesDelete(svc.context, id)
	res, err := svc.client.RegistriesApi.RegistriesDeleteExecute(req)
	return res, err
}

func (svc *registriesService) Patch(id string, input containerregistry.PatchRegistryInput) (
	containerregistry.RegistryResponse, *containerregistry.APIResponse, error,
) {
	req := svc.client.RegistriesApi.RegistriesPatch(svc.context, id).PatchRegistryInput(input)
	reg, res, err := svc.client.RegistriesApi.RegistriesPatchExecute(req)
	return reg, res, err
}

func (svc *registriesService) Put(id string, input containerregistry.PutRegistryInput) (
	containerregistry.PutRegistryOutput, *containerregistry.APIResponse, error,
) {
	req := svc.client.RegistriesApi.RegistriesPut(svc.context, id).PutRegistryInput(input)
	reg, res, err := svc.client.RegistriesApi.RegistriesPutExecute(req)
	return reg, res, err
}
