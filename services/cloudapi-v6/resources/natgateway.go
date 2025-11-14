package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

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
	List(datacenterId string) (NatGateways, *Response, error)
	Get(datacenterId, natGatewayId string) (*NatGateway, *Response, error)
	Create(datacenterId string, input NatGateway) (*NatGateway, *Response, error)
	Update(datacenterId, natGatewayId string, input NatGatewayProperties) (*NatGateway, *Response, error)
	Delete(datacenterId, natGatewayId string) (*Response, error)
	ListRules(datacenterId, natGatewayId string) (NatGatewayRules, *Response, error)
	GetRule(datacenterId, natGatewayId, ruleId string) (*NatGatewayRule, *Response, error)
	CreateRule(datacenterId, natGatewayId string, input NatGatewayRule) (*NatGatewayRule, *Response, error)
	UpdateRule(datacenterId, natGatewayId, ruleId string, input NatGatewayRuleProperties) (*NatGatewayRule, *Response, error)
	DeleteRule(datacenterId, natGatewayId, ruleId string) (*Response, error)
	ListFlowLogs(datacenterId, natGatewayId string) (FlowLogs, *Response, error)
	GetFlowLog(datacenterId, natGatewayId, flowlogId string) (*FlowLog, *Response, error)
	CreateFlowLog(datacenterId, natGatewayId string, input FlowLog) (*FlowLog, *Response, error)
	UpdateFlowLog(datacenterId, natGatewayId, flowlogId string, input *FlowLogProperties) (*FlowLog, *Response, error)
	DeleteFlowLog(datacenterId, natGatewayId, flowlogId string) (*Response, error)
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

func (ds *natGatewaysService) List(datacenterId string) (NatGateways, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysGet(ds.context, datacenterId)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysGetExecute(req)
	return NatGateways{s}, &Response{*res}, err
}

func (ds *natGatewaysService) Get(datacenterId, natGatewayId string) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayId(ds.context, datacenterId, natGatewayId)
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFindByNatGatewayIdExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Create(datacenterId string, input NatGateway) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysPost(ds.context, datacenterId).NatGateway(input.NatGateway)
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysPostExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Update(datacenterId, natGatewayId string, input NatGatewayProperties) (*NatGateway, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysPatch(ds.context, datacenterId, natGatewayId).NatGatewayProperties(input.NatGatewayProperties)
	datacenter, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysPatchExecute(req)
	return &NatGateway{datacenter}, &Response{*res}, err
}

func (ds *natGatewaysService) Delete(datacenterId, natGatewayId string) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysDelete(ds.context, datacenterId, natGatewayId)
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysDeleteExecute(req)
	return &Response{*res}, err
}

func (ds *natGatewaysService) ListRules(datacenterId, natGatewayId string) (NatGatewayRules, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesGet(ds.context, datacenterId, natGatewayId)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesGetExecute(req)
	return NatGatewayRules{s}, &Response{*res}, err
}

func (ds *natGatewaysService) GetRule(datacenterId, natGatewayId, ruleId string) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ds.context, datacenterId, natGatewayId, ruleId)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleIdExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) CreateRule(datacenterId, natGatewayId string, input NatGatewayRule) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPost(ds.context, datacenterId, natGatewayId).NatGatewayRule(input.NatGatewayRule)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPostExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) UpdateRule(datacenterId, natGatewayId, ruleId string, input NatGatewayRuleProperties) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatch(ds.context, datacenterId, natGatewayId, ruleId).NatGatewayRuleProperties(input.NatGatewayRuleProperties)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatchExecute(req)
	return &NatGatewayRule{s}, &Response{*res}, err
}

func (ds *natGatewaysService) DeleteRule(datacenterId, natGatewayId, ruleId string) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ds.context, datacenterId, natGatewayId, ruleId)
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDeleteExecute(req)
	return &Response{*res}, err
}

func (ds *natGatewaysService) ListFlowLogs(datacenterId, natGatewayId string) (FlowLogs, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsGet(ds.context, datacenterId, natGatewayId)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsGetExecute(req)
	return FlowLogs{s}, &Response{*res}, err
}

func (ds *natGatewaysService) GetFlowLog(datacenterId, natGatewayId, flowlogId string) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsFindByFlowLogId(ds.context, datacenterId, natGatewayId, flowlogId)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsFindByFlowLogIdExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) CreateFlowLog(datacenterId, natGatewayId string, input FlowLog) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPost(ds.context, datacenterId, natGatewayId).NatGatewayFlowLog(input.FlowLog)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPostExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) UpdateFlowLog(datacenterId, natGatewayId, flowlogId string, input *FlowLogProperties) (*FlowLog, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPatch(ds.context, datacenterId, natGatewayId, flowlogId).NatGatewayFlowLogProperties(input.FlowLogProperties)
	s, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsPatchExecute(req)
	return &FlowLog{s}, &Response{*res}, err
}

func (ds *natGatewaysService) DeleteFlowLog(datacenterId, natGatewayId, flowlogId string) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsDelete(ds.context, datacenterId, natGatewayId, flowlogId)
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysFlowlogsDeleteExecute(req)
	return &Response{*res}, err
}
