package resources

import (
	"context"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	List(datacenterId, serverId, nicId string, param ListQueryParams) (FirewallRules, *Response, error)
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

func (svc *firewallRulesService) List(datacenterId, serverId, nicId string, params ListQueryParams) (FirewallRules, *Response, error) {
	req := svc.client.NicApi.DatacentersServersNicsFirewallrulesGet(svc.context, datacenterId, serverId, nicId)
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
	}
	rules, resp, err := svc.client.NicApi.DatacentersServersNicsFirewallrulesGetExecute(req)
	return FirewallRules{rules}, &Response{*resp}, err
}

func (svc *firewallRulesService) Get(datacenterId, serverId, nicId, firewallRuleId string) (*FirewallRule, *Response, error) {
	req := svc.client.NicApi.DatacentersServersNicsFirewallrulesFindById(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	rule, resp, err := svc.client.NicApi.DatacentersServersNicsFirewallrulesFindByIdExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Create(datacenterId, serverId, nicId string, input FirewallRule) (*FirewallRule, *Response, error) {
	req := svc.client.NicApi.DatacentersServersNicsFirewallrulesPost(svc.context, datacenterId, serverId, nicId).Firewallrule(input.FirewallRule)
	rule, resp, err := svc.client.NicApi.DatacentersServersNicsFirewallrulesPostExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Update(datacenterId, serverId, nicId, firewallRuleId string, input FirewallRuleProperties) (*FirewallRule, *Response, error) {
	req := svc.client.NicApi.DatacentersServersNicsFirewallrulesPatch(svc.context, datacenterId, serverId, nicId, firewallRuleId).Firewallrule(input.FirewallruleProperties)
	rule, resp, err := svc.client.NicApi.DatacentersServersNicsFirewallrulesPatchExecute(req)
	return &FirewallRule{rule}, &Response{*resp}, err
}

func (svc *firewallRulesService) Delete(datacenterId, serverId, nicId, firewallRuleId string) (*Response, error) {
	req := svc.client.NicApi.DatacentersServersNicsFirewallrulesDelete(svc.context, datacenterId, serverId, nicId, firewallRuleId)
	_, resp, err := svc.client.NicApi.DatacentersServersNicsFirewallrulesDeleteExecute(req)
	return &Response{*resp}, err
}
