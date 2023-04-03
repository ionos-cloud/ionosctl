package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
)

type Nic struct {
	ionoscloud.Nic
}

type NicProperties struct {
	ionoscloud.NicProperties
}

type Nics struct {
	ionoscloud.Nics
}

type LanNics struct {
	ionoscloud.LanNics
}

type BalancedNics struct {
	ionoscloud.BalancedNics
}

// NicsService is a wrapper around ionoscloud.Nic
type NicsService interface {
	List(datacenterId, serverId string, params ListQueryParams) (Nics, *Response, error)
	Get(datacenterId, serverId, nicId string, params QueryParams) (*Nic, *Response, error)
	Create(datacenterId, serverId string, input Nic, params QueryParams) (*Nic, *Response, error)
	Update(datacenterId, serverId, nicId string, input NicProperties, params QueryParams) (*Nic, *Response, error)
	Delete(datacenterId, serverId, nicId string, params QueryParams) (*Response, error)
}

type nicsService struct {
	client  *ionoscloud.APIClient
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
