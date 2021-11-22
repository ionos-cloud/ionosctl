package resources

import (
	"context"
	"github.com/fatih/structs"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	Get(datacenterId, serverId, nicId string) (*Nic, *Response, error)
	Create(datacenterId, serverId, name string, ips []string, dhcp bool, lan int32) (*Nic, *Response, error)
	Update(datacenterId, serverId, nicId string, input NicProperties) (*Nic, *Response, error)
	Delete(datacenterId, serverId, nicId string) (*Response, error)
}

type nicsService struct {
	client  *Client
	context context.Context
}

var _ NicsService = &nicsService{}

func NewNicService(client *Client, ctx context.Context) NicsService {
	return &nicsService{
		client:  client,
		context: ctx,
	}
}

func (ns *nicsService) List(datacenterId, serverId string, params ListQueryParams) (Nics, *Response, error) {
	req := ns.client.NicApi.DatacentersServersNicsGet(ns.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
	}
	nics, resp, err := ns.client.NicApi.DatacentersServersNicsGetExecute(req)
	return Nics{nics}, &Response{*resp}, err
}

func (ns *nicsService) Get(datacenterId, serverId, nicId string) (*Nic, *Response, error) {
	req := ns.client.NicApi.DatacentersServersNicsFindById(ns.context, datacenterId, serverId, nicId)
	nic, resp, err := ns.client.NicApi.DatacentersServersNicsFindByIdExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Create(datacenterId, serverId, name string, ips []string, dhcp bool, lan int32) (*Nic, *Response, error) {
	input := ionoscloud.Nic{
		Properties: &ionoscloud.NicProperties{
			Name: &name,
			Ips:  &ips,
			Dhcp: &dhcp,
			Lan:  &lan,
		},
	}
	req := ns.client.NicApi.DatacentersServersNicsPost(ns.context, datacenterId, serverId).Nic(input)
	nic, resp, err := ns.client.NicApi.DatacentersServersNicsPostExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Update(datacenterId, serverId, nicId string, input NicProperties) (*Nic, *Response, error) {
	req := ns.client.NicApi.DatacentersServersNicsPatch(ns.context, datacenterId, serverId, nicId).Nic(input.NicProperties)
	nic, resp, err := ns.client.NicApi.DatacentersServersNicsPatchExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) Delete(datacenterId, serverId, nicId string) (*Response, error) {
	req := ns.client.NicApi.DatacentersServersNicsDelete(ns.context, datacenterId, serverId, nicId)
	_, resp, err := ns.client.NicApi.DatacentersServersNicsDeleteExecute(req)
	return &Response{*resp}, err
}
