package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	client  *ionoscloud.APIClient
	context context.Context
}

var _ ContractsService = &contractsService{}

func NewContractService(client *client.Client, ctx context.Context) ContractsService {
	return &contractsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *contractsService) Get() (Contracts, *Response, error) {
	req := s.client.ContractResourcesApi.ContractsGet(s.context)
	req = client.ApplyFilters(req, viper.GetStringSlice(constants.FlagFilters))
	contracts, resp, err := s.client.ContractResourcesApi.ContractsGetExecute(req)
	return Contracts{contracts}, &Response{*resp}, err
}
