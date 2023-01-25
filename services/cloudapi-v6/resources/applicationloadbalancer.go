package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	"github.com/fatih/structs"

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
	Get(datacenterId, applicationLoadBalancerId string, params QueryParams) (*ApplicationLoadBalancer, *Response, error)
	Create(datacenterId string, input ApplicationLoadBalancer, params QueryParams) (*ApplicationLoadBalancer, *Response, error)
	Update(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerProperties, params QueryParams) (*ApplicationLoadBalancer, *Response, error)
	Delete(datacenterId, applicationLoadBalancerId string, params QueryParams) (*Response, error)
	ListForwardingRules(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (ApplicationLoadBalancerForwardingRules, *Response, error)
	GetForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	CreateForwardingRule(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerForwardingRule, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	UpdateForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, input *ApplicationLoadBalancerForwardingRuleProperties, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error)
	DeleteForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, params QueryParams) (*Response, error)
	ListFlowLogs(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (FlowLogs, *Response, error)
	GetFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, params QueryParams) (*FlowLog, *Response, error)
	CreateFlowLog(datacenterId, applicationLoadBalancerId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error)
	UpdateFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, input *FlowLogProperties, params QueryParams) (*FlowLog, *Response, error)
	DeleteFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, params QueryParams) (*Response, error)
}

type applicationLoadBalancersService struct {
	client  *config.Client
	context context.Context
}

var _ ApplicationLoadBalancersService = &applicationLoadBalancersService{}

func NewApplicationLoadBalancerService(client *config.Client, ctx context.Context) ApplicationLoadBalancersService {
	return &applicationLoadBalancersService{
		client:  client,
		context: ctx,
	}
}

func (svc *applicationLoadBalancersService) List(datacenterId string, params ListQueryParams) (ApplicationLoadBalancers, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGet(svc.context, datacenterId)
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
	applicationLoadBalancers, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersGetExecute(req)
	return ApplicationLoadBalancers{applicationLoadBalancers}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Get(datacenterId, applicationLoadBalancerId string, params QueryParams) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerId(svc.context, datacenterId, applicationLoadBalancerId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFindByApplicationLoadBalancerIdExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Create(datacenterId string, input ApplicationLoadBalancer, params QueryParams) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPost(svc.context, datacenterId).ApplicationLoadBalancer(input.ApplicationLoadBalancer)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPostExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Update(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerProperties, params QueryParams) (*ApplicationLoadBalancer, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatch(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerProperties(input.ApplicationLoadBalancerProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	applicationLoadBalancer, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersPatchExecute(req)
	return &ApplicationLoadBalancer{applicationLoadBalancer}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) Delete(datacenterId, applicationLoadBalancerId string, params QueryParams) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDelete(svc.context, datacenterId, applicationLoadBalancerId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) ListForwardingRules(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (ApplicationLoadBalancerForwardingRules, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGet(svc.context, datacenterId, applicationLoadBalancerId)
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
	applicationLoadBalancerRules, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesGetExecute(req)
	return ApplicationLoadBalancerForwardingRules{applicationLoadBalancerRules}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) GetForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleId(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesFindByForwardingRuleIdExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) CreateForwardingRule(datacenterId, applicationLoadBalancerId string, input ApplicationLoadBalancerForwardingRule, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPost(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerForwardingRule(input.ApplicationLoadBalancerForwardingRule)
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPostExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) UpdateForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, input *ApplicationLoadBalancerForwardingRuleProperties, params QueryParams) (*ApplicationLoadBalancerForwardingRule, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatch(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId).ApplicationLoadBalancerForwardingRuleProperties(input.ApplicationLoadBalancerForwardingRuleProperties)
	applicationLoadBalancerRule, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesPatchExecute(req)
	return &ApplicationLoadBalancerForwardingRule{applicationLoadBalancerRule}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) DeleteForwardingRule(datacenterId, applicationLoadBalancerId, forwardingRuleId string, params QueryParams) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDelete(svc.context, datacenterId, applicationLoadBalancerId, forwardingRuleId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersForwardingrulesDeleteExecute(req)
	return &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) ListFlowLogs(datacenterId, applicationLoadBalancerId string, params ListQueryParams) (FlowLogs, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsGet(svc.context, datacenterId, applicationLoadBalancerId)
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
	flowLogs, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsGetExecute(req)
	return FlowLogs{flowLogs}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) GetFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsFindByFlowLogId(svc.context, datacenterId, applicationLoadBalancerId, flowLogId)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsFindByFlowLogIdExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) CreateFlowLog(datacenterId, applicationLoadBalancerId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPost(svc.context, datacenterId, applicationLoadBalancerId).ApplicationLoadBalancerFlowLog(input.FlowLog)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPostExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) UpdateFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, input *FlowLogProperties, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPatch(svc.context, datacenterId, applicationLoadBalancerId, flowLogId).ApplicationLoadBalancerFlowLogProperties(input.FlowLogProperties)
	flowLog, resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsPatchExecute(req)
	return &FlowLog{flowLog}, &Response{*resp}, err
}

func (svc *applicationLoadBalancersService) DeleteFlowLog(datacenterId, applicationLoadBalancerId, flowLogId string, params QueryParams) (*Response, error) {
	req := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsDelete(svc.context, datacenterId, applicationLoadBalancerId, flowLogId)
	resp, err := svc.client.ApplicationLoadBalancersApi.DatacentersApplicationloadbalancersFlowlogsDeleteExecute(req)
	return &Response{*resp}, err
}
