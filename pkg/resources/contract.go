package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Contract struct {
	ionoscloud.Contract
}

// ContractsService is a wrapper around ionoscloud.Contract
type ContractsService interface {
	Get() (Contract, *Response, error)
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

func (s *contractsService) Get() (Contract, *Response, error) {
	req := s.client.ContractApi.ContractsGet(s.context)
	contract, resp, err := s.client.ContractApi.ContractsGetExecute(req)
	return Contract{contract}, &Response{*resp}, err
}
