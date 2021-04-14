package resources

import (
	"context"

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
	List(datacenterId, serverId string) (Nics, *Response, error)
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

func (ns *nicsService) List(datacenterId, serverId string) (Nics, *Response, error) {
	req := ns.client.NicApi.DatacentersServersNicsGet(ns.context, datacenterId, serverId)
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

func (ns *nicsService) AttachToLan(datacenterId, lanId, nicId string) (*Nic, *Response, error) {
	req := ns.client.LanApi.DatacentersLansNicsPost(ns.context, datacenterId, lanId)
	req = req.Nic(ionoscloud.Nic{Id: &nicId})
	nic, resp, err := ns.client.LanApi.DatacentersLansNicsPostExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *nicsService) ListAttachedToLan(datacenterId, lanId string) (LanNics, *Response, error) {
	req := ns.client.LanApi.DatacentersLansNicsGet(ns.context, datacenterId, lanId)
	nics, resp, err := ns.client.LanApi.DatacentersLansNicsGetExecute(req)
	return LanNics{nics}, &Response{*resp}, err
}

func (ns *nicsService) GetAttachedToLan(datacenterId, lanId, nicId string) (*Nic, *Response, error) {
	req := ns.client.LanApi.DatacentersLansNicsFindById(ns.context, datacenterId, lanId, nicId)
	n, resp, err := ns.client.LanApi.DatacentersLansNicsFindByIdExecute(req)
	return &Nic{n}, &Response{*resp}, err
}
