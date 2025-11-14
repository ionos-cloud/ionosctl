package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

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
	List(datacenterId, serverId, nicId string, params ListQueryParams) (FirewallRules, *Response, error)
	Get(datacenterId, serverId, nicId, firewallRuleId string, params QueryParams) (*FirewallRule, *Response, error)
	Create(datacenterId, serverId, nicId string, input FirewallRule, params QueryParams) (*FirewallRule, *Response, error)
	Update(datacenterId, serverId, nicId, firewallRuleId string, input FirewallRuleProperties, params QueryParams) (*FirewallRule, *Response, error)
	Delete(datacenterId, serverId, nicId, firewallRuleId string, params QueryParams) (*Response, error)
}

type firewallRulesService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ FirewallRulesService = &firewallRulesService{}

func NewFirewallRuleService(client *client.Client, ctx context.Context) FirewallRulesService {
	return &firewallRulesService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (svc *firewallRulesService) List(datacenterId, serverId, nicId string, params ListQueryParams) (FirewallRules, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGet(svc.context, datacenterId, serverId, nicId)
	rules, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesGetExecute(req)
	return FirewallRules{rules}, &Response{*resp}, err
}

func (svc *firewallRulesService) Get(datacenterId, serverId, nicId, firewallRuleId string, params QueryParams) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindByIdExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Create(datacenterId, serverId, nicId string, input FirewallRule, params QueryParams) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPost(svc.context, datacenterId, serverId, nicId).Firewallrule(input.FirewallRule)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPostExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Update(datacenterId, serverId, nicId, firewallRuleId string, input FirewallRuleProperties, params QueryParams) (*FirewallRule, *Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatch(svc.context, datacenterId, serverId, nicId, firewallRuleId).Firewallrule(input.FirewallruleProperties)
	rule, resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesPatchExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Delete(datacenterId, serverId, nicId, firewallRuleId string, params QueryParams) (*Response, error) {
	req := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDelete(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	resp, err := svc.client.FirewallRulesApi.DatacentersServersNicsFirewallrulesDeleteExecute(req)
	return &Response{*resp}, err
}
