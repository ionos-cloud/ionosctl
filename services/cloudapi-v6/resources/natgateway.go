package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type NatGateway struct {
	ionoscloud.NatGateway
}

type NatGatewayRule struct {
	ionoscloud.NatGatewayRule
}

type NatGatewayLanProperties struct {
	ionoscloud.NatGatewayLanProperties
}

type NatGatewayRuleProperties struct {
	ionoscloud.NatGatewayRuleProperties
}

type NatGatewayRules struct {
	ionoscloud.NatGatewayRules
}

type NatGatewayProperties struct {
	ionoscloud.NatGatewayProperties
}

type NatGateways struct {
	ionoscloud.NatGateways
}

// NatGatewaysService is a wrapper around ionoscloud.NatGateway
type NatGatewaysService interface {
	List(datacenterId string, params ListQueryParams) (NatGateways, *Response, error)
	Get(datacenterId, natGatewayId string, params QueryParams) (*NatGateway, *Response, error)
	Create(datacenterId string, input NatGateway, params QueryParams) (*NatGateway, *Response, error)
	Update(datacenterId, natGatewayId string, input NatGatewayProperties, params QueryParams) (*NatGateway, *Response, error)
	Delete(datacenterId, natGatewayId string, params QueryParams) (*Response, error)
	ListRules(datacenterId, natGatewayId string, params ListQueryParams) (NatGatewayRules, *Response, error)
	GetRule(datacenterId, natGatewayId, ruleId string, params QueryParams) (*NatGatewayRule, *Response, error)
	CreateRule(datacenterId, natGatewayId string, input NatGatewayRule, params QueryParams) (*NatGatewayRule, *Response, error)
	UpdateRule(datacenterId, natGatewayId, ruleId string, input NatGatewayRuleProperties, params QueryParams) (*NatGatewayRule, *Response, error)
	DeleteRule(datacenterId, natGatewayId, ruleId string, params QueryParams) (*Response, error)
	ListFlowLogs(datacenterId, natGatewayId string, params ListQueryParams) (FlowLogs, *Response, error)
	GetFlowLog(datacenterId, natGatewayId, flowlogId string, params QueryParams) (*FlowLog, *Response, error)
	CreateFlowLog(datacenterId, natGatewayId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error)
	UpdateFlowLog(datacenterId, natGatewayId, flowlogId string, input *FlowLogProperties, params QueryParams) (*FlowLog, *Response, error)
	DeleteFlowLog(datacenterId, natGatewayId, flowlogId string, params QueryParams) (*Response, error)
}

type natGatewaysService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ NatGatewaysService = &natGatewaysService{}

func NewNatGatewayService(client *client.Client, ctx context.Context) NatGatewaysService {
	return &natGatewaysService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ds *natGatewaysService) List(datacenterId string, params ListQueryParams) (NatGateways, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysGet(ds.context, datacenterId)
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
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysGetExecute(req)
	return NatGateways{s}, &Response{*res}, err
}

func (ds *natGatewaysService) Get(datacenterId, natGatewayId string, params QueryParams) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ds.context, datacenterId, natGatewayId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayIdExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Create(datacenterId string, input NatGateway, params QueryParams) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysPost(ds.context, datacenterId).NatGateway(input.NatGateway)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysPostExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Update(datacenterId, natGatewayId string, input NatGatewayProperties, params QueryParams) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysPatch(ds.context, datacenterId, natGatewayId).NatGatewayProperties(input.NatGatewayProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysPatchExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Delete(datacenterId, natGatewayId string, params QueryParams) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysDelete(ds.context, datacenterId, natGatewayId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysDeleteExecute(req)
	return &Response{*res}, err
}

func (ds *natGatewaysService) ListRules(datacenterId, natGatewayId string, params ListQueryParams) (NatGatewayRules, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ds.context, datacenterId, natGatewayId)
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
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesGetExecute(req)
	return NatGatewayRules{s}, &Response{*res}, err
}

func (ds *natGatewaysService) GetRule(datacenterId, natGatewayId, ruleId string, params QueryParams) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ds.context, datacenterId, natGatewayId, ruleId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleIdExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) CreateRule(datacenterId, natGatewayId string, input NatGatewayRule, params QueryParams) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPost(ds.context, datacenterId, natGatewayId).NatGatewayRule(input.NatGatewayRule)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPostExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) UpdateRule(datacenterId, natGatewayId, ruleId string, input NatGatewayRuleProperties, params QueryParams) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatch(ds.context, datacenterId, natGatewayId, ruleId).NatGatewayRuleProperties(input.NatGatewayRuleProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatchExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) DeleteRule(datacenterId, natGatewayId, ruleId string, params QueryParams) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ds.context, datacenterId, natGatewayId, ruleId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDeleteExecute(req)
	return &Response{*res}, err
}

func (ds *natGatewaysService) ListFlowLogs(datacenterId, natGatewayId string, params ListQueryParams) (FlowLogs, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsGet(ds.context, datacenterId, natGatewayId)
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
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsGetExecute(req)
	return FlowLogs{s}, &Response{*res}, err
}

func (ds *natGatewaysService) GetFlowLog(datacenterId, natGatewayId, flowlogId string, params QueryParams) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsFindByFlowLogId(ds.context, datacenterId, natGatewayId, flowlogId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsFindByFlowLogIdExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) CreateFlowLog(datacenterId, natGatewayId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPost(ds.context, datacenterId, natGatewayId).NatGatewayFlowLog(input.FlowLog)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPostExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) UpdateFlowLog(datacenterId, natGatewayId, flowlogId string, input *FlowLogProperties, params QueryParams) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPatch(ds.context, datacenterId, natGatewayId, flowlogId).NatGatewayFlowLogProperties(input.FlowLogProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPatchExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) DeleteFlowLog(datacenterId, natGatewayId, flowlogId string, params QueryParams) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsDelete(ds.context, datacenterId, natGatewayId, flowlogId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsDeleteExecute(req)
	return &Response{*res}, err
}
