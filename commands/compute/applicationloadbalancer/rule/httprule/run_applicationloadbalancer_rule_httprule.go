package httprule

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunApplicationLoadBalancerRuleHttpRule(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId,
		cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgName, constants.FlagType)
}

func PreRunApplicationLoadBalancerRuleHttpRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgName},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgAll},
	)
}

func RunAlbRuleHttpRuleList(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	c.Verbose("Getting HttpRules")

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	properties, ok := ng.GetPropertiesOk()
	if !ok || properties == nil {
		return errors.New("error getting rule properties")
	}

	httpRules, ok := properties.GetHttpRulesOk()
	if !ok || httpRules == nil {
		return errors.New("error getting rule http rules")
	}

	return c.Printer(allAlbRuleHttpRuleCols).Print(*httpRules)
}

func RunAlbRuleHttpRuleAdd(c *core.CommandConfig) error {
	var httpRuleItems []ionoscloud.ApplicationLoadBalancerHttpRule

	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
	c.Verbose("Getting HttpRules from ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	ngOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}

	if properties, ok := ngOld.GetPropertiesOk(); ok && properties != nil {
		if httpRulesOk, ok := properties.GetHttpRulesOk(); ok && httpRulesOk != nil {
			httpRuleItems = *httpRulesOk
		}
	}

	httpRuleNew := getRuleHttpRuleInfo(c)
	c.Verbose("Adding the new HttpRule to the existing HttpRules")

	httpRuleItems = append(httpRuleItems, httpRuleNew.ApplicationLoadBalancerHttpRule)
	c.Verbose("Updating ForwardingRule with the new HttpRules")

	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		&resources.ApplicationLoadBalancerForwardingRuleProperties{
			ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
				HttpRules: &httpRuleItems,
			},
		},
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allAlbRuleHttpRuleCols).Print(httpRuleNew.ApplicationLoadBalancerHttpRule)
}

func RunAlbRuleHttpRuleRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
		c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

		resp, err := RemoveAllHTTPRules(c)
		if err != nil {
			return err
		}
		if resp != nil {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}

		return nil
	}

	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove forwarding rule http rule", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Getting HttpRules")

	frOld, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return err
	}

	c.Verbose("Removing the HTTP Rule from the existing HTTP Rules")

	proper, err := getRuleHttpRulesRemove(c, frOld)
	if err != nil {
		return err
	}

	c.Verbose("Updating ForwardingRule with the new HTTP Rules")

	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		proper,
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Application Load Balancer HTTP Rule successfully deleted")

	return nil

}

func RemoveAllHTTPRules(c *core.CommandConfig) (*resources.Response, error) {
	c.Msg("Forwarding Rule HTTP Rules to be deleted:")

	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if err != nil {
		return nil, err
	}

	propertiesOk, ok := applicationLoadBalancerRules.GetPropertiesOk()
	if !ok || propertiesOk == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer Forwarding Rule properties")
	}

	httpRulesOk, ok := propertiesOk.GetHttpRulesOk()
	if !ok || httpRulesOk == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer HTTP Rules")
	}

	if len(*httpRulesOk) <= 0 {
		return nil, fmt.Errorf("no Application Load Balancer HTTP Rules found")
	}

	for _, httpRuleOk := range *httpRulesOk {
		if nameOk, ok := httpRuleOk.GetNameOk(); ok && nameOk != nil {
			c.Msg("Forwarding Rule HTTP Rule Name: %v", *nameOk)
		}

		if typeOk, ok := httpRuleOk.GetTypeOk(); ok && typeOk != nil {
			c.Msg("Forwarding Rule HTTP Rule Type: %v", *typeOk)
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Forwarding Rule HTTP Rules", viper.GetBool(constants.ArgForce)) {
		return nil, fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Deleting all the Forwarding Rule HTTP Rules...")

	propertiesOk.SetHttpRules([]ionoscloud.ApplicationLoadBalancerHttpRule{})
	_, resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		&resources.ApplicationLoadBalancerForwardingRuleProperties{
			ApplicationLoadBalancerForwardingRuleProperties: *propertiesOk,
		},
	)
	if resp != nil {
		c.Verbose("Request Id: %v", request.GetId(resp))
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return nil, err
	}

	c.Msg("Application Load Balancer HTTP Rules successfully deleted")
	return resp, err
}

func getRuleHttpRuleInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRule {
	// Set Application Load Balancer HTTP Rule Properties
	httprule := resources.ApplicationLoadBalancerHttpRule{}

	httprule.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))

	httprule.SetType(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))
	c.Verbose("Property Type set: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "FORWARD") {
		httprule.SetTargetGroup(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		c.Verbose("Property TargetGroup set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	}

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "REDIRECT") {
		httprule.SetLocation(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
		c.Verbose("Property Location set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))

		httprule.SetDropQuery(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgQuery)))
		c.Verbose("Property DropQuery set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgQuery)))

		httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
		c.Verbose("Property StatusCode set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
	}

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)), "STATIC") {
		httprule.SetResponseMessage(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMessage)))
		c.Verbose("Property ResponseMessage set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMessage)))

		httprule.SetContentType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))
		c.Verbose("Property ContentType set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgContentType)))

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)) {
			httprule.SetStatusCode(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
			c.Verbose("Property StatusCode set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgStatusCode)))
		} else {
			httprule.SetStatusCode(503)
			c.Verbose("Property StatusCode set with the default value: %v", 503)
		}
	}

	httpRuleCondition := getRuleHttpRuleConditionInfo(c)
	httprule.SetConditions([]ionoscloud.ApplicationLoadBalancerHttpRuleCondition{
		httpRuleCondition.ApplicationLoadBalancerHttpRuleCondition,
	})

	c.Verbose("Setting Condition to HttpRule")

	return httprule
}

func getRuleHttpRuleConditionInfo(c *core.CommandConfig) resources.ApplicationLoadBalancerHttpRuleCondition {
	// Set Application Load Balancer HTTP Rule Condition Properties
	httpRuleCondition := resources.ApplicationLoadBalancerHttpRuleCondition{}
	httpRuleCondition.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)))

	c.Verbose("Property Condition Type set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)))

	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "SOURCE_IP") {
		httpRuleCondition.SetCondition(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
		c.Verbose("Property Condition set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)))
	}

	httpRuleCondition.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))

	c.Verbose("Property Condition Negate set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))

	if strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "COOKIES") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "HEADER") ||
		strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionType)), "QUERY") {
		httpRuleCondition.SetKey(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
		c.Verbose("Property Condition Key set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionKey)))
	}

	if !strings.EqualFold(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCondition)), "EXISTS") {
		httpRuleCondition.SetValue(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
		c.Verbose("Property Condition Value set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgConditionValue)))
	}

	return httpRuleCondition
}

func getRuleHttpRulesRemove(c *core.CommandConfig, frOld *resources.ApplicationLoadBalancerForwardingRule) (*resources.ApplicationLoadBalancerForwardingRuleProperties, error) {
	httpRuleItems := make([]ionoscloud.ApplicationLoadBalancerHttpRule, 0)

	properties, ok := frOld.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not get Application Load Balancer Forwarding Rule properties")
	}

	c.Verbose("Getting Properties from the Forwarding Rule")

	httpRules, ok := properties.GetHttpRulesOk()
	if !ok || httpRules == nil {
		return nil, fmt.Errorf("coudl not get Application Load Balancer HTTP Rules")
	}

	c.Verbose("Getting HTTP Rules from the Forwarding Rule Properties")

	for _, httpRuleItem := range *httpRules {
		removeName := false

		if nameOk, ok := httpRuleItem.GetNameOk(); ok && nameOk != nil {
			if *nameOk == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
				removeName = true
				c.Verbose("Found HTTP Rule with Name: %v", *nameOk)
			}
		}

		// If the Http rule with the unique name is found, continue.
		// If not, add it to the Forwarding Rule properties.
		if removeName {
			continue
		} else {
			httpRuleItems = append(httpRuleItems, httpRuleItem)
		}
	}

	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{
			HttpRules: &httpRuleItems,
		},
	}, nil
}

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}

func PreRunDcApplicationLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId)
}
