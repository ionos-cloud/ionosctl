package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// TokenService is a wrapper around ionoscloud.Registry
type TokenService interface {
	Get(id string, registryId string) (sdkgo.TokenResponse, *shared.APIResponse, error)
	List(registryId string) (sdkgo.TokensResponse, *shared.APIResponse, error)
	Post(input sdkgo.PostTokenInput, registryId string) (sdkgo.PostTokenOutput, *shared.APIResponse, error)
	Delete(id string, registryId string) (*shared.APIResponse, error)
	Patch(id string, input sdkgo.PatchTokenInput, registryId string) (sdkgo.TokenResponse, *shared.APIResponse, error)
	Put(id string, input sdkgo.PutTokenInput, registryId string) (sdkgo.PutTokenOutput, *shared.APIResponse, error)
}

type tokenService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ TokenService = &tokenService{}

func NewTokenService(client *client2.Client, ctx context.Context) TokenService {
	return &tokenService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *tokenService) List(registryId string) (sdkgo.TokensResponse, *shared.APIResponse, error) {
	req := svc.client.TokensApi.RegistriesTokensGet(svc.context, registryId)
	tokenList, res, err := svc.client.TokensApi.RegistriesTokensGetExecute(req)
	return tokenList, res, err
}

func (svc *tokenService) Post(input sdkgo.PostTokenInput, registryId string) (sdkgo.PostTokenOutput, *shared.APIResponse, error) {
	req := svc.client.TokensApi.RegistriesTokensPost(svc.context, registryId).PostTokenInput(input)
	tokenPost, res, err := svc.client.TokensApi.RegistriesTokensPostExecute(req)
	return tokenPost, res, err
}

func (svc *tokenService) Get(id string, registryId string) (sdkgo.TokenResponse, *shared.APIResponse, error) {
	req := svc.client.TokensApi.RegistriesTokensFindById(svc.context, registryId, id)
	reg, res, err := svc.client.TokensApi.RegistriesTokensFindByIdExecute(req)
	return reg, res, err
}

func (svc *tokenService) Delete(id string, registryId string) (*shared.APIResponse, error) {
	req := svc.client.TokensApi.RegistriesTokensDelete(svc.context, registryId, id)
	res, err := svc.client.TokensApi.RegistriesTokensDeleteExecute(req)
	return res, err
}

func (svc *tokenService) Patch(id string, input sdkgo.PatchTokenInput, registryId string) (
	sdkgo.TokenResponse, *shared.APIResponse, error,
) {
	req := svc.client.TokensApi.RegistriesTokensPatch(svc.context, registryId, id).PatchTokenInput(input)
	reg, res, err := svc.client.TokensApi.RegistriesTokensPatchExecute(req)
	return reg, res, err
}

func (svc *tokenService) Put(id string, input sdkgo.PutTokenInput, registryId string) (
	sdkgo.PutTokenOutput, *shared.APIResponse, error,
) {
	req := svc.client.TokensApi.RegistriesTokensPut(svc.context, registryId, id).PutTokenInput(input)
	reg, res, err := svc.client.TokensApi.RegistriesTokensPutExecute(req)
	return reg, res, err
}
