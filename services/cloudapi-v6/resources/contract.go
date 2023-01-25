package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	"github.com/fatih/structs"

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
	Get(params QueryParams) (Contracts, *Response, error)
}

type contractsService struct {
	client  *config.Client
	context context.Context
}

var _ ContractsService = &contractsService{}

func NewContractService(client *config.Client, ctx context.Context) ContractsService {
	return &contractsService{
		client:  client,
		context: ctx,
	}
}

func (s *contractsService) Get(params QueryParams) (Contracts, *Response, error) {
	req := s.client.ContractResourcesApi.ContractsGet(s.context)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	contracts, resp, err := s.client.ContractResourcesApi.ContractsGetExecute(req)
	return Contracts{contracts}, &Response{*resp}, err
}
