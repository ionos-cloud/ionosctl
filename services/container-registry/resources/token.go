package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

// TokenService is a wrapper around ionoscloud.Registry
type TokenService interface {
	Get(id string, registryId string) (sdkgo.TokenResponse, *Response, error)
	List(registryId string) (sdkgo.TokensResponse, *Response, error)
	Post(input sdkgo.PostTokenInput, registryId string) (sdkgo.PostTokenOutput, *Response, error)
	Delete(id string, registryId string) (*Response, error)
	Patch(id string, input sdkgo.PatchTokenInput, registryId string) (sdkgo.TokenResponse, *Response, error)
	Put(id string, input sdkgo.PutTokenInput, registryId string) (sdkgo.PutTokenOutput, *Response, error)
}

type tokenService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ TokenService = &tokenService{}

func NewTokenService(client *config.Client, ctx context.Context) TokenService {
	return &tokenService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *tokenService) List(registryId string) (sdkgo.TokensResponse, *Response, error) {
	req := svc.client.TokensApi.RegistriesTokensGet(svc.context, registryId)
	tokenList, res, err := svc.client.TokensApi.RegistriesTokensGetExecute(req)
	return tokenList, &Response{*res}, err
}

func (svc *tokenService) Post(input sdkgo.PostTokenInput, registryId string) (sdkgo.PostTokenOutput, *Response, error) {
	req := svc.client.TokensApi.RegistriesTokensPost(svc.context, registryId).PostTokenInput(input)
	tokenPost, res, err := svc.client.TokensApi.RegistriesTokensPostExecute(req)
	return tokenPost, &Response{*res}, err
}

func (svc *tokenService) Get(id string, registryId string) (sdkgo.TokenResponse, *Response, error) {
	req := svc.client.TokensApi.RegistriesTokensFindById(svc.context, registryId, id)
	reg, res, err := svc.client.TokensApi.RegistriesTokensFindByIdExecute(req)
	return reg, &Response{*res}, err
}

func (svc *tokenService) Delete(id string, registryId string) (*Response, error) {
	req := svc.client.TokensApi.RegistriesTokensDelete(svc.context, registryId, id)
	res, err := svc.client.TokensApi.RegistriesTokensDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *tokenService) Patch(id string, input sdkgo.PatchTokenInput, registryId string) (
	sdkgo.TokenResponse, *Response, error,
) {
	req := svc.client.TokensApi.RegistriesTokensPatch(svc.context, registryId, id).PatchTokenInput(input)
	reg, res, err := svc.client.TokensApi.RegistriesTokensPatchExecute(req)
	return reg, &Response{*res}, err
}

func (svc *tokenService) Put(id string, input sdkgo.PutTokenInput, registryId string) (
	sdkgo.PutTokenOutput, *Response, error,
) {
	req := svc.client.TokensApi.RegistriesTokensPut(svc.context, registryId, id).PutTokenInput(input)
	reg, res, err := svc.client.TokensApi.RegistriesTokensPutExecute(req)
	return reg, &Response{*res}, err
}
