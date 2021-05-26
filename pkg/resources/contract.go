package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Contract struct {
	ionoscloud.Contract
}

type Contracts struct {
	ionoscloud.Contracts
}

// ContractsService is a wrapper around ionoscloud.Contract
type ContractsService interface {
	Get() (Contracts, *Response, error)
}

type contractsService struct {
	client  *Client
	context context.Context
}

var _ ContractsService = &contractsService{}

func NewContractService(client *Client, ctx context.Context) ContractsService {
	return &contractsService{
		client:  client,
		context: ctx,
	}
}

func (s *contractsService) Get() (Contracts, *Response, error) {
	req := s.client.ContractResourcesApi.ContractsGet(s.context)
	contracts, resp, err := s.client.ContractResourcesApi.ContractsGetExecute(req)
	return Contracts{contracts}, &Response{*resp}, err
}
