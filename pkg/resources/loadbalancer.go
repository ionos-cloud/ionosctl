package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
}

type loadbalancersService struct {
	client  *Client
	context context.Context
}

var _ LoadbalancersService = &loadbalancersService{}

func NewLoadbalancerService(client *Client, ctx context.Context) LoadbalancersService {
	return &loadbalancersService{
		client:  client,
		context: ctx,
	}
}

func (ls *loadbalancersService) List(datacenterId string) (Loadbalancers, *Response, error) {
	req := ls.client.LoadBalancerApi.DatacentersLoadbalancersGet(ls.context, datacenterId)
	s, res, err := ls.client.LoadBalancerApi.DatacentersLoadbalancersGetExecute(req)
	return Loadbalancers{s}, &Response{*res}, err
}

func (ls *loadbalancersService) Get(datacenterId, loadbalancerId string) (*Loadbalancer, *Response, error) {
	req := ls.client.LoadBalancerApi.DatacentersLoadbalancersFindById(ls.context, datacenterId, loadbalancerId)
	s, res, err := ls.client.LoadBalancerApi.DatacentersLoadbalancersFindByIdExecute(req)
	return &Loadbalancer{s}, &Response{*res}, err
}

func (ls *loadbalancersService) Create(datacenterId, name string, dhcp bool) (*Loadbalancer, *Response, error) {
	s := ionoscloud.Loadbalancer{
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &name,
			Dhcp: &dhcp,
		},
	}
	req := ls.client.LoadBalancerApi.DatacentersLoadbalancersPost(ls.context, datacenterId).Loadbalancer(s)
	loadbalancer, res, err := ls.client.LoadBalancerApi.DatacentersLoadbalancersPostExecute(req)
	return &Loadbalancer{loadbalancer}, &Response{*res}, err
}

func (ls *loadbalancersService) Update(datacenterId, loadbalancerId string, input LoadbalancerProperties) (*Loadbalancer, *Response, error) {
	req := ls.client.LoadBalancerApi.DatacentersLoadbalancersPatch(ls.context, datacenterId, loadbalancerId).Loadbalancer(input.LoadbalancerProperties)
	loadbalancer, resp, err := ls.client.LoadBalancerApi.DatacentersLoadbalancersPatchExecute(req)
	return &Loadbalancer{loadbalancer}, &Response{*resp}, err
}

func (ls *loadbalancersService) Delete(datacenterId, loadbalancerId string) (*Response, error) {
	req := ls.client.LoadBalancerApi.DatacentersLoadbalancersDelete(ls.context, datacenterId, loadbalancerId)
	_, res, err := ls.client.LoadBalancerApi.DatacentersLoadbalancersDeleteExecute(req)
	return &Response{*res}, err
}
