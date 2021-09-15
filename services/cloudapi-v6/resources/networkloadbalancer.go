package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type NetworkLoadBalancer struct {
	ionoscloud.NetworkLoadBalancer
}

type NetworkLoadBalancerPut struct {
	ionoscloud.NetworkLoadBalancerPut
}

type NetworkLoadBalancerProperties struct {
	ionoscloud.NetworkLoadBalancerProperties
}

type NetworkLoadBalancers struct {
	ionoscloud.NetworkLoadBalancers
}

type NetworkLoadBalancerForwardingRule struct {
	ionoscloud.NetworkLoadBalancerForwardingRule
}

type NetworkLoadBalancerForwardingRuleProperties struct {
	ionoscloud.NetworkLoadBalancerForwardingRuleProperties
}

type NetworkLoadBalancerForwardingRuleHealthCheck struct {
	ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck
}

type NetworkLoadBalancerForwardingRuleTarget struct {
	ionoscloud.NetworkLoadBalancerForwardingRuleTarget
}

type NetworkLoadBalancerForwardingRuleTargetHealthCheck struct {
	ionoscloud.NetworkLoadBalancerForwardingRuleTargetHealthCheck
}

type NetworkLoadBalancerForwardingRules struct {
	ionoscloud.NetworkLoadBalancerForwardingRules
}

// NetworkLoadBalancersService is a wrapper around ionoscloud.NetworkLoadBalancer
type NetworkLoadBalancersService interface {
	List(datacenterId string) (NetworkLoadBalancers, *Response, error)
	Get(datacenterId, networkLoadBalancerId string) (*NetworkLoadBalancer, *Response, error)
	Create(datacenterId string, input NetworkLoadBalancer) (*NetworkLoadBalancer, *Response, error)
	Update(datacenterId, networkLoadBalancerId string, input NetworkLoadBalancerProperties) (*NetworkLoadBalancer, *Response, error)
	Delete(datacenterId, networkLoadBalancerId string) (*Response, error)
	ListForwardingRules(datacenterId, networkLoadBalancerId string) (NetworkLoadBalancerForwardingRules, *Response, error)
	GetForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*NetworkLoadBalancerForwardingRule, *Response, error)
	CreateForwardingRule(datacenterId, networkLoadBalancerId string, input NetworkLoadBalancerForwardingRule) (*NetworkLoadBalancerForwardingRule, *Response, error)
	UpdateForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string, input *NetworkLoadBalancerForwardingRuleProperties) (*NetworkLoadBalancerForwardingRule, *Response, error)
	DeleteForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*Response, error)
	ListFlowLogs(datacenterId, networkLoadBalancerId string) (FlowLogs, *Response, error)
	GetFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*FlowLog, *Response, error)
	CreateFlowLog(datacenterId, networkLoadBalancerId string, input FlowLog) (*FlowLog, *Response, error)
	UpdateFlowLog(datacenterId, networkLoadBalancerId, flowLogId string, input *FlowLogProperties) (*FlowLog, *Response, error)
	DeleteFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*Response, error)
}

type networkLoadBalancersService struct {
	client  *Client
	context context.Context
}

var _ NetworkLoadBalancersService = &networkLoadBalancersService{}

func NewNetworkLoadBalancerService(client *Client, ctx context.Context) NetworkLoadBalancersService {
	return &networkLoadBalancersService{
		client:  client,
		context: ctx,
	}
}

func (svc *networkLoadBalancersService) List(datacenterId string) (NetworkLoadBalancers, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGet(svc.context, datacenterId)
	networkLoadBalancers, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersGetExecute(req)
	return NetworkLoadBalancers{networkLoadBalancers}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) Get(datacenterId, networkLoadBalancerId string) (*NetworkLoadBalancer, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerId(svc.context, datacenterId, networkLoadBalancerId)
	networkLoadBalancer, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFindByNetworkLoadBalancerIdExecute(req)
	return &NetworkLoadBalancer{networkLoadBalancer}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) Create(datacenterId string, input NetworkLoadBalancer) (*NetworkLoadBalancer, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPost(svc.context, datacenterId).NetworkLoadBalancer(input.NetworkLoadBalancer)
	networkLoadBalancer, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPostExecute(req)
	return &NetworkLoadBalancer{networkLoadBalancer}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) Update(datacenterId, networkLoadBalancerId string, input NetworkLoadBalancerProperties) (*NetworkLoadBalancer, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPatch(svc.context, datacenterId, networkLoadBalancerId).NetworkLoadBalancerProperties(input.NetworkLoadBalancerProperties)
	networkLoadBalancer, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersPatchExecute(req)
	return &NetworkLoadBalancer{networkLoadBalancer}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) Delete(datacenterId, networkLoadBalancerId string) (*Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersDelete(svc.context, datacenterId, networkLoadBalancerId)
	resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *networkLoadBalancersService) ListForwardingRules(datacenterId, networkLoadBalancerId string) (NetworkLoadBalancerForwardingRules, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGet(svc.context, datacenterId, networkLoadBalancerId)
	networkLoadBalancerRules, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesGetExecute(req)
	return NetworkLoadBalancerForwardingRules{networkLoadBalancerRules}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) GetForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*NetworkLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleId(svc.context, datacenterId, networkLoadBalancerId, forwardingRuleId)
	networkLoadBalancerRule, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesFindByForwardingRuleIdExecute(req)
	return &NetworkLoadBalancerForwardingRule{networkLoadBalancerRule}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) CreateForwardingRule(datacenterId, networkLoadBalancerId string, input NetworkLoadBalancerForwardingRule) (*NetworkLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPost(svc.context, datacenterId, networkLoadBalancerId).NetworkLoadBalancerForwardingRule(input.NetworkLoadBalancerForwardingRule)
	networkLoadBalancerRule, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPostExecute(req)
	return &NetworkLoadBalancerForwardingRule{networkLoadBalancerRule}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) UpdateForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string, input *NetworkLoadBalancerForwardingRuleProperties) (*NetworkLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPatch(svc.context, datacenterId, networkLoadBalancerId, forwardingRuleId).NetworkLoadBalancerForwardingRuleProperties(input.NetworkLoadBalancerForwardingRuleProperties)
	networkLoadBalancerRule, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesPatchExecute(req)
	return &NetworkLoadBalancerForwardingRule{networkLoadBalancerRule}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) DeleteForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDelete(svc.context, datacenterId, networkLoadBalancerId, forwardingRuleId)
	resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersForwardingrulesDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *networkLoadBalancersService) ListFlowLogs(datacenterId, networkLoadBalancerId string) (FlowLogs, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsGet(svc.context, datacenterId, networkLoadBalancerId)
	flowLogs, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsGetExecute(req)
	return FlowLogs{flowLogs}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) GetFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*FlowLog, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsFindByFlowLogId(svc.context, datacenterId, networkLoadBalancerId, flowLogId)
	flowLog, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsFindByFlowLogIdExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) CreateFlowLog(datacenterId, networkLoadBalancerId string, input FlowLog) (*FlowLog, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPost(svc.context, datacenterId, networkLoadBalancerId).NetworkLoadBalancerFlowLog(input.FlowLog)
	flowLog, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPostExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) UpdateFlowLog(datacenterId, networkLoadBalancerId, flowLogId string, input *FlowLogProperties) (*FlowLog, *Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPatch(svc.context, datacenterId, networkLoadBalancerId, flowLogId).NetworkLoadBalancerFlowLogProperties(input.FlowLogProperties)
	flowLog, resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsPatchExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *networkLoadBalancersService) DeleteFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*Response, error) {
	req := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsDelete(svc.context, datacenterId, networkLoadBalancerId, flowLogId)
	resp, err := svc.client.NetworkLoadBalancersApi.DatacentersNetworkloadbalancersFlowlogsDeleteExecute(req)
	return &Response{*resp}, err
}
