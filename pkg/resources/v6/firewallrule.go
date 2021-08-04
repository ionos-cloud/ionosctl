package v6

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type FirewallRule struct {
	ionoscloud.FirewallRule
}

type FirewallRuleProperties struct {
	ionoscloud.FirewallruleProperties
}

type FirewallRules struct {
	ionoscloud.FirewallRules
}

// FirewallRulesService is a wrapper around ionoscloud.FirewallRule
type FirewallRulesService interface {
	List(datacenterId, serverId, nicId string) (FirewallRules, *Response, error)
	Get(datacenterId, serverId, nicId, firewallRuleId string) (*FirewallRule, *Response, error)
	Create(datacenterId, serverId, nicId string, input FirewallRule) (*FirewallRule, *Response, error)
	Update(datacenterId, serverId, nicId, firewallRuleId string, input FirewallRuleProperties) (*FirewallRule, *Response, error)
	Delete(datacenterId, serverId, nicId, firewallRuleId string) (*Response, error)
}

type firewallRulesService struct {
	client  *Client
	context context.Context
}

var _ FirewallRulesService = &firewallRulesService{}

func NewFirewallRuleService(client *Client, ctx context.Context) FirewallRulesService {
	return &firewallRulesService{
		client:  client,
		context: ctx,
	}
}

func (svc *firewallRulesService) List(datacenterId, serverId, nicId string) (FirewallRules, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(svc.context, datacenterId, serverId, nicId)
	rules, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGetExecute(req)
	return FirewallRules{rules}, &Response{*resp}, err
}

func (svc *firewallRulesService) Get(datacenterId, serverId, nicId, firewallRuleId string) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindByIdExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Create(datacenterId, serverId, nicId string, input FirewallRule) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(svc.context, datacenterId, serverId, nicId).Firewallrule(input.FirewallRule)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPostExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Update(datacenterId, serverId, nicId, firewallRuleId string, input FirewallRuleProperties) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(svc.context, datacenterId, serverId, nicId, firewallRuleId).Firewallrule(input.FirewallruleProperties)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatchExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Delete(datacenterId, serverId, nicId, firewallRuleId string) (*Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDelete(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDeleteExecute(req)
	return &Response{*resp}, err
}
