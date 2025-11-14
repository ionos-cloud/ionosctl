package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type ApplicationLoadBalancer struct {
	ionoscloud.ApplicationLoadBalancer
}

type ApplicationLoadBalancerPut struct {
	ionoscloud.ApplicationLoadBalancerPut
}

type ApplicationLoadBalancerProperties struct {
	ionoscloud.ApplicationLoadBalancerProperties
}

type ApplicationLoadBalancers struct {
	ionoscloud.ApplicationLoadBalancers
}

type ApplicationLoadBalancerForwardingRule struct {
	ionoscloud.ApplicationLoadBalancerForwardingRule
}

type ApplicationLoadBalancerHttpRule struct {
	ionoscloud.ApplicationLoadBalancerHttpRule
}

type ApplicationLoadBalancerHttpRuleCondition struct {
	ionoscloud.ApplicationLoadBalancerHttpRuleCondition
}

type ApplicationLoadBalancerForwardingRuleProperties struct {
	ionoscloud.ApplicationLoadBalancerForwardingRuleProperties
}

type ApplicationLoadBalancerForwardingRules struct {
	ionoscloud.ApplicationLoadBalancerForwardingRules
}

// ApplicationLoadBalancersService is a wrapper around ionoscloud.ApplicationLoadBalancer
type ApplicationLoadBalancersService interface {
	List(datacenterId string, params ListQueryParams) (ApplicationLoadBalancers, *Response, error)
	Get(datacenterId, applicationLoadBalancerId string) (*ApplicationLoadBalancer, *Response, error)
	Create(datacenterId string, input ApplicationLoadBalancer) (*ApplicationLoadBalancer, *Response, error)
	Update(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerProperties) (*ApplicationLoadBalancer, *Response, error)
	Delete(datacenterId, applicationLoadBalancerId string) (*Response, error)
	ListForwardingRules(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (ApplicationLoadBalancerForwardingRules, *Response, error)
	GetForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	CreateForwardingRule(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerForwardingRule) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	UpdateForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, input *ApplicationLoadBalancerForwardingRuleProperties) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	DeleteForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string) (*Response, error)
	ListFlowLogs(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (FlowLogs, *Response, error)
	GetFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string) (*FlowLog, *Response, error)
	CreateFlowLog(datacenterId, applicationLoadBalancerId string, input FlowLog) (*FlowLog, *Response, error)
	UpdateFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, input *FlowLogProperties) (*FlowLog, *Response, error)
	DeleteFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string) (*Response, error)
}

type applicationLoadBalancersService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ ApplicationLoadBalancersService = &applicationLoadBalancersService{}

func NewApplicationLoadBalancerService(client *client.Client, ctx context.Context) ApplicationLoadBalancersService {
	return &applicationLoadBalancersService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (svc *applicationLoadBalancersService) List(datacenterId string, params ListQueryParams) (ApplicationLoadBalancers, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(svc.context, datacenterId)
	applicationLoadBalancers, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGetExecute(req)
	return ApplicationLoadBalancers{applicationLoadBalancers}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Get(datacenterId, applicationLoadBalancerId string) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(svc.context, datacenterId, applicationLoadBalancerId)
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerIdExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Create(datacenterId string, input ApplicationLoadBalancer) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPost(svc.context, datacenterId).ApplicationLoadBalancer(input.ApplicationLoadBalancer)
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPostExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Update(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerProperties) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatch(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerProperties(input.ApplicationLoadBalancerProperties)
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatchExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Delete(datacenterId, applicationLoadBalancerId string) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDelete(svc.context, datacenterId, applicationLoadBalancerId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) ListForwardingRules(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (ApplicationLoadBalancerForwardingRules, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGet(svc.context, datacenterId, applicationLoadBalancerId)
	applicationLoadBalancerRules, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGetExecute(req)
	return ApplicationLoadBalancerForwardingRules{applicationLoadBalancerRules}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) GetForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId)
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleIdExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) CreateForwardingRule(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerForwardingRule) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPost(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerForwardingRule(input.ApplicationLoadBalancerForwardingRule)
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPostExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) UpdateForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, input *ApplicationLoadBalancerForwardingRuleProperties) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatch(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId).ApplicationLoadBalancerForwardingRuleProperties(input.ApplicationLoadBalancerForwardingRuleProperties)
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatchExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) DeleteForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDelete(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) ListFlowLogs(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (FlowLogs, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsGet(svc.context, datacenterId, applicationLoadBalancerId)
	flowLogs, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsGetExecute(req)
	return FlowLogs{flowLogs}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) GetFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsFindByFlowLogId(svc.context, datacenterId, applicationLoadBalancerId, flowLogId)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsFindByFlowLogIdExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) CreateFlowLog(datacenterId, applicationLoadBalancerId string, input FlowLog) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPost(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerFlowLog(input.FlowLog)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPostExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) UpdateFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, input *FlowLogProperties) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPatch(svc.context, datacenterId, applicationLoadBalancerId, flowLogId).ApplicationLoadBalancerFlowLogProperties(input.FlowLogProperties)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPatchExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) DeleteFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsDelete(svc.context, datacenterId, applicationLoadBalancerId, flowLogId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsDeleteExecute(req)
	return &Response{*resp}, err
}
