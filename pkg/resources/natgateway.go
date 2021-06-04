package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type NatGateway struct {
	ionoscloud.NatGateway
}

type NatGatewayRule struct {
	ionoscloud.NatGatewayRule
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
}

type natGatewaysService struct {
	client  *Client
	context context.Context
}

var _ NatGatewaysService = &natGatewaysService{}

func NewNatGatewayService(client *Client, ctx context.Context) NatGatewaysService {
	return &natGatewaysService{
		client:  client,
		context: ctx,
	}
}

func (ds *natGatewaysService) List(datacenterId string) (NatGateways, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysGet(ds.context, datacenterId)
	dcs, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysGetExecute(req)
	return NatGateways{dcs}, &Response{*res}, err
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
	dcs, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesGetExecute(req)
	return NatGatewayRules{dcs}, &Response{*res}, err
}

func (ds *natGatewaysService) GetRule(datacenterId, natGatewayId, ruleId string) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleId(ds.context, datacenterId, natGatewayId, ruleId)
	dcs, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesFindByNatGatewayRuleIdExecute(req)
	return &NatGatewayRule{dcs}, &Response{*res}, err
}

func (ds *natGatewaysService) CreateRule(datacenterId, natGatewayId string, input NatGatewayRule) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPost(ds.context, datacenterId, natGatewayId).NatGatewayRule(input.NatGatewayRule)
	dcs, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPostExecute(req)
	return &NatGatewayRule{dcs}, &Response{*res}, err
}

func (ds *natGatewaysService) UpdateRule(datacenterId, natGatewayId, ruleId string, input NatGatewayRuleProperties) (*NatGatewayRule, *Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatch(ds.context, datacenterId, natGatewayId, ruleId).NatGatewayRuleProperties(input.NatGatewayRuleProperties)
	dcs, res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesPatchExecute(req)
	return &NatGatewayRule{dcs}, &Response{*res}, err
}

func (ds *natGatewaysService) DeleteRule(datacenterId, natGatewayId, ruleId string) (*Response, error) {
	req := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDelete(ds.context, datacenterId, natGatewayId, ruleId)
	res, err := ds.client.NATGatewaysApi.DatacentersNatgatewaysRulesDeleteExecute(req)
	return &Response{*res}, err
}
