package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/viper"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Loadbalancer struct {
	ionoscloud.Loadbalancer
}

type LoadbalancerProperties struct {
	ionoscloud.LoadbalancerProperties
}

type Loadbalancers struct {
	ionoscloud.Loadbalancers
}

// LoadbalancersService is a wrapper around ionoscloud.Loadbalancer
type LoadbalancersService interface {
	List(datacenterId string) (Loadbalancers, *Response, error)
	Get(datacenterId, loadbalancerId string) (*Loadbalancer, *Response, error)
	Create(datacenterId, name string, dhcp bool) (*Loadbalancer, *Response, error)
	Update(datacenterId, loadbalancerId string, input LoadbalancerProperties) (*Loadbalancer, *Response, error)
	Delete(datacenterId, loadbalancerId string) (*Response, error)
	AttachNic(datacenterId, loadbalancerId, nicId string) (*Nic, *Response, error)
	ListNics(datacenterId, loadbalancerId string) (BalancedNics, *Response, error)
	GetNic(datacenterId, loadbalancerId, nicId string) (*Nic, *Response, error)
	DetachNic(datacenterId, loadbalancerId, nicId string) (*Response, error)
}

type loadbalancersService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ LoadbalancersService = &loadbalancersService{}

func NewLoadbalancerService(client *client.Client, ctx context.Context) LoadbalancersService {
	return &loadbalancersService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ls *loadbalancersService) List(datacenterId string) (Loadbalancers, *Response, error) {
	req := ls.client.LoadBalancersApi.DatacentersLoadbalancersGet(ls.context, datacenterId)
	req = client.ApplyFilters(req, viper.GetStringSlice(constants.FlagFilters))
	s, res, err := ls.client.LoadBalancersApi.DatacentersLoadbalancersGetExecute(req)
	return Loadbalancers{s}, &Response{*res}, err
}

func (ls *loadbalancersService) Get(datacenterId, loadbalancerId string) (*Loadbalancer, *Response, error) {
	req := ls.client.LoadBalancersApi.DatacentersLoadbalancersFindById(ls.context, datacenterId, loadbalancerId)
	s, res, err := ls.client.LoadBalancersApi.DatacentersLoadbalancersFindByIdExecute(req)
	return &Loadbalancer{s}, &Response{*res}, err
}

func (ls *loadbalancersService) Create(datacenterId, name string, dhcp bool) (*Loadbalancer, *Response, error) {
	s := ionoscloud.Loadbalancer{
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &name,
			Dhcp: &dhcp,
		},
	}
	req := ls.client.LoadBalancersApi.DatacentersLoadbalancersPost(ls.context, datacenterId).Loadbalancer(s)
	loadbalancer, res, err := ls.client.LoadBalancersApi.DatacentersLoadbalancersPostExecute(req)
	return &Loadbalancer{loadbalancer}, &Response{*res}, err
}

func (ls *loadbalancersService) Update(datacenterId, loadbalancerId string, input LoadbalancerProperties) (*Loadbalancer, *Response, error) {
	req := ls.client.LoadBalancersApi.DatacentersLoadbalancersPatch(ls.context, datacenterId, loadbalancerId).Loadbalancer(input.LoadbalancerProperties)
	loadbalancer, resp, err := ls.client.LoadBalancersApi.DatacentersLoadbalancersPatchExecute(req)
	return &Loadbalancer{loadbalancer}, &Response{*resp}, err
}

func (ls *loadbalancersService) Delete(datacenterId, loadbalancerId string) (*Response, error) {
	req := ls.client.LoadBalancersApi.DatacentersLoadbalancersDelete(ls.context, datacenterId, loadbalancerId)
	res, err := ls.client.LoadBalancersApi.DatacentersLoadbalancersDeleteExecute(req)
	return &Response{*res}, err
}

func (ns *loadbalancersService) AttachNic(datacenterId, loadbalancerId, nicId string) (*Nic, *Response, error) {
	input := ionoscloud.Nic{Id: &nicId}
	req := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsPost(ns.context, datacenterId, loadbalancerId).Nic(input)
	nic, resp, err := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsPostExecute(req)
	return &Nic{nic}, &Response{*resp}, err
}

func (ns *loadbalancersService) ListNics(datacenterId, loadbalancerId string) (BalancedNics, *Response, error) {
	req := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsGet(ns.context, datacenterId, loadbalancerId)
	req = client.ApplyFilters(req, viper.GetStringSlice(constants.FlagFilters))
	nics, resp, err := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsGetExecute(req)
	return BalancedNics{nics}, &Response{*resp}, err
}

func (ns *loadbalancersService) GetNic(datacenterId, loadbalancerId, nicId string) (*Nic, *Response, error) {
	req := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsFindByNicId(ns.context, datacenterId, loadbalancerId, nicId)
	n, resp, err := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsFindByNicIdExecute(req)
	return &Nic{n}, &Response{*resp}, err
}

func (ns *loadbalancersService) DetachNic(datacenterId, loadbalancerId, nicId string) (*Response, error) {
	req := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsDelete(ns.context, datacenterId, loadbalancerId, nicId)
	resp, err := ns.client.LoadBalancersApi.DatacentersLoadbalancersBalancednicsDeleteExecute(req)
	return &Response{*resp}, err
}
