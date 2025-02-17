package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"

	"github.com/fatih/structs"
)

type IpBlock struct {
	compute.IpBlock
}

type IpBlockProperties struct {
	compute.IpBlockProperties
}

type IpBlocks struct {
	compute.IpBlocks
}

type IpConsumer struct {
	compute.IpConsumer
}

// IpBlocksService is a wrapper around compute.IpBlock
type IpBlocksService interface {
	List(params ListQueryParams) (IpBlocks, *Response, error)
	Get(IpBlockId string, params QueryParams) (*IpBlock, *Response, error)
	Create(name, location string, size int32, params QueryParams) (*IpBlock, *Response, error)
	Update(ipBlockId string, input IpBlockProperties, params QueryParams) (*IpBlock, *Response, error)
	Delete(ipBlockId string, params QueryParams) (*Response, error)
}

type ipBlocksService struct {
	client  *compute.APIClient
	context context.Context
}

var _ IpBlocksService = &ipBlocksService{}

func NewIpBlockService(client *client.Client, ctx context.Context) IpBlocksService {
	return &ipBlocksService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (svc *ipBlocksService) List(params ListQueryParams) (IpBlocks, *Response, error) {
	req := svc.client.IPBlocksApi.IpblocksGet(svc.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	s, res, err := svc.client.IPBlocksApi.IpblocksGetExecute(req)
	return IpBlocks{s}, &Response{*res}, err
}

func (svc *ipBlocksService) Get(ipBlockId string, params QueryParams) (*IpBlock, *Response, error) {
	req := svc.client.IPBlocksApi.IpblocksFindById(svc.context, ipBlockId)
	s, res, err := svc.client.IPBlocksApi.IpblocksFindByIdExecute(req)
	return &IpBlock{s}, &Response{*res}, err
}

func (svc *ipBlocksService) Create(name, location string, size int32, params QueryParams) (*IpBlock, *Response, error) {
	i := compute.IpBlock{
		Properties: &compute.IpBlockProperties{
			Location: &location,
			Size:     &size,
		},
	}
	if name != "" {
		i.Properties.SetName(name)
	}
	req := svc.client.IPBlocksApi.IpblocksPost(svc.context).Ipblock(i)
	ipBlock, res, err := svc.client.IPBlocksApi.IpblocksPostExecute(req)
	return &IpBlock{ipBlock}, &Response{*res}, err
}

func (svc *ipBlocksService) Update(ipBlockId string, input IpBlockProperties, params QueryParams) (*IpBlock, *Response, error) {
	req := svc.client.IPBlocksApi.IpblocksPatch(svc.context, ipBlockId).Ipblock(input.IpBlockProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	ipBlock, resp, err := svc.client.IPBlocksApi.IpblocksPatchExecute(req)
	return &IpBlock{ipBlock}, &Response{*resp}, err
}

func (svc *ipBlocksService) Delete(ipBlockId string, params QueryParams) (*Response, error) {
	req := svc.client.IPBlocksApi.IpblocksDelete(svc.context, ipBlockId)
	res, err := svc.client.IPBlocksApi.IpblocksDeleteExecute(req)
	return &Response{*res}, err
}
