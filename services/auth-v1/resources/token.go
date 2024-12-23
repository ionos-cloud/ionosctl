package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
)

type Token struct {
	sdkgoauth.Token
}

type Jwt struct {
	sdkgoauth.Jwt
}

type Tokens struct {
	sdkgoauth.Tokens
}

type Response struct {
	sdkgoauth.APIResponse
}

type DeleteResponse struct {
	sdkgoauth.DeleteResponse
}

// TokensService is a wrapper around ionoscloud.Token
type TokensService interface {
	List(contractNumber int32) (Tokens, *Response, error)
	Get(tokenId string, contractNumber int32) (*Token, *Response, error)
	Create(contractNumber int32, ttl int32) (*Jwt, *Response, error)
	DeleteByID(tokenId string, contractNumber int32) (*DeleteResponse, *Response, error)
	DeleteByCriteria(criteria string, contractNumber int32) (*DeleteResponse, *Response, error)
}

type tokensService struct {
	client  *sdkgoauth.APIClient
	context context.Context
}

var _ TokensService = &tokensService{}

func NewTokenService(client *client.Client, ctx context.Context) TokensService {
	return &tokensService{
		client:  client.AuthClient,
		context: ctx,
	}
}

func (ts *tokensService) List(contractNumber int32) (Tokens, *Response, error) {
	req := ts.client.TokensApi.TokensGet(ts.context)
	if contractNumber != 0 {
		req = req.XContractNumber(contractNumber)
	}
	dcs, res, err := ts.client.TokensApi.TokensGetExecute(req)
	return Tokens{dcs}, &Response{*res}, err
}

func (ts *tokensService) Get(tokenId string, contractNumber int32) (*Token, *Response, error) {
	req := ts.client.TokensApi.TokensFindById(ts.context, tokenId)
	if contractNumber != 0 {
		req = req.XContractNumber(contractNumber)
	}
	token, res, err := ts.client.TokensApi.TokensFindByIdExecute(req)
	return &Token{token}, &Response{*res}, err
}

func (ts *tokensService) Create(contractNumber int32, ttl int32) (*Jwt, *Response, error) {
	req := ts.client.TokensApi.TokensGenerate(ts.context)
	if contractNumber != 0 {
		req = req.XContractNumber(contractNumber)
	}
	if ttl != 0 {
		req = req.Ttl(ttl)
	}
	token, res, err := ts.client.TokensApi.TokensGenerateExecute(req)
	return &Jwt{token}, &Response{*res}, err
}

func (ts *tokensService) DeleteByID(tokenId string, contractNumber int32) (*DeleteResponse, *Response, error) {
	req := ts.client.TokensApi.TokensDeleteById(ts.context, tokenId)
	if contractNumber != 0 {
		req = req.XContractNumber(contractNumber)
	}
	tokenDeleteById, res, err := ts.client.TokensApi.TokensDeleteByIdExecute(req)
	return &DeleteResponse{tokenDeleteById}, &Response{*res}, err
}

// DeleteByCriteria removes all tokens based on criteria: EXPIRED, CURRENT or ALL
func (ts *tokensService) DeleteByCriteria(criteria string, contractNumber int32) (*DeleteResponse, *Response, error) {
	req := ts.client.TokensApi.TokensDeleteByCriteria(ts.context).Criteria(criteria)
	if contractNumber != 0 {
		req = req.XContractNumber(contractNumber)
	}
	tokenDeleteByCriteria, res, err := ts.client.TokensApi.TokensDeleteByCriteriaExecute(req)
	return &DeleteResponse{tokenDeleteByCriteria}, &Response{*res}, err
}
