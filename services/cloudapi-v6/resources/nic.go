package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"

	"github.com/fatih/structs"
)

type Nic struct {
	compute.Nic
}

type NicProperties struct {
	compute.NicProperties
}

type Nics struct {
	compute.Nics
}

type LanNics struct {
	compute.LanNics
}

type BalancedNics struct {
	compute.BalancedNics
}

// NicsService is a wrapper around compute.Nic
type NicsService interface {
	List(datacenterId, serverId string, params ListQueryParams) (Nics, *Response, error)
	Get(datacenterId, serverId, nicId string, params QueryParams) (*Nic, *Response, error)
	Create(datacenterId, serverId string, input Nic, params QueryParams) (*Nic, *Response, error)
	Update(datacenterId, serverId, nicId string, input NicProperties, params QueryParams) (*Nic, *Response, error)
	Delete(datacenterId, serverId, nicId string, params QueryParams) (*Response, error)
}

type nicsService struct {
	client  *compute.APIClient
	context context.Context
}

var _ NicsService = &nicsService{}

func NewNicService(client *client.Client, ctx context.Context) NicsService {
	return &nicsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ns *nicsService) List(datacenterId, serverId string, params ListQueryParams) (Nics, *Response, error) {
	req := ns.client.NetworkInterfacesApi.DatacentersServersNicsGet(ns.context, datacenterId, serverId)
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
	nics, resp, err := ns.client.NetworkInterfacesApi.DatacentersServersNicsGetExecute(req)
	return Nics{nics}, &Response{*resp}, err
}

func (ns *nicsService) Get(datacenterId, serverId, nicId string, params QueryParams) (*Nic, *Response, error) {
	req := ns.client.NetworkInterfacesApi.DatacentersServersNicsFindById(ns.context, datacenterId, serverId, nicId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	nic, resp, err := ns.client.NetworkInterfacesApi.DatacentersServersNicsFindByIdExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Create(datacenterId, serverId string, input Nic, params QueryParams) (*Nic, *Response, error) {
	req := ns.client.NetworkInterfacesApi.DatacentersServersNicsPost(ns.context, datacenterId, serverId).Nic(input.Nic)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	nic, resp, err := ns.client.NetworkInterfacesApi.DatacentersServersNicsPostExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Update(datacenterId, serverId, nicId string, input NicProperties, params QueryParams) (*Nic, *Response, error) {
	req := ns.client.NetworkInterfacesApi.DatacentersServersNicsPatch(ns.context, datacenterId, serverId, nicId).Nic(input.NicProperties)
	nic, resp, err := ns.client.NetworkInterfacesApi.DatacentersServersNicsPatchExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Delete(datacenterId, serverId, nicId string, params QueryParams) (*Response, error) {
	req := ns.client.NetworkInterfacesApi.DatacentersServersNicsDelete(ns.context, datacenterId, serverId, nicId)
	resp, err := ns.client.NetworkInterfacesApi.DatacentersServersNicsDeleteExecute(req)
	return &Response{*resp}, err
}
