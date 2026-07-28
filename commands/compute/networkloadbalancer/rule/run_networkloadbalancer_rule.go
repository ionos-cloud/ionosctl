package rule

import (
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

func PreRunNetworkLoadBalancerRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId); err != nil {
		return err
	}
	return nil
}

func PreRunNetworkLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgListenerIp, cloudapiv6.ArgListenerPort)
}

func PreRunDcNetworkLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId)
}

func PreRunDcNetworkLoadBalancerForwardingRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func RunNetworkLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(nlbForwardingRules.NetworkLoadBalancerForwardingRules)
}

func RunNetworkLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	c.Verbose("Network Load Balancer Forwarding Rule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancerForwardingRule)
}

func RunNetworkLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	proper := getForwardingRulePropertiesSet(c)

	if !proper.HasProtocol() {
		proper.SetProtocol(string(ionoscloud.TCP))
	}

	if !proper.HasAlgorithm() {
		proper.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	}

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		proper.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		resources.NetworkLoadBalancerForwardingRule{
			NetworkLoadBalancerForwardingRule: ionoscloud.NetworkLoadBalancerForwardingRule{
				Properties: &proper.NetworkLoadBalancerForwardingRuleProperties,
			},
		},
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancerForwardingRule)
}

func RunNetworkLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	input := getForwardingRulePropertiesSet(c)
	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		input.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancerForwardingRule)
}

func RunNetworkLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNetworkLoadBalancerForwardingRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete network load balancer forwarding rule", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Network Load Balancer Forwarding Rule with id: %v...", ruleId)

	resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, ruleId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Network Load Balancer Forwarding Rule successfully deleted")
	return nil
}

func getForwardingRulePropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleProperties {
	input := ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)) {
		algorithm := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
		input.SetAlgorithm(algorithm)

		c.Verbose("Property Algorithm set: %v", algorithm)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)) {
		listenerIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp))
		input.SetListenerIp(listenerIp)

		c.Verbose("Property ListenerIp set: %v", listenerIp)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)) {
		listenerPort := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort))
		input.SetListenerPort(listenerPort)

		c.Verbose("Property ListenerPort set: %v", listenerPort)
	}

	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: input,
	}
}

func getForwardingRuleHealthCheckPropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleHealthCheck {
	inputHealth := ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealth.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)) {
		inputHealth.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)) {
		inputHealth.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)) {
		inputHealth.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
	}

	return &resources.NetworkLoadBalancerForwardingRuleHealthCheck{
		NetworkLoadBalancerForwardingRuleHealthCheck: inputHealth,
	}
}

func DeleteAllNetworkLoadBalancerForwardingRules(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Network Load Balancer ID: %v", loadBalancerId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.NetworkLoadBalancerForwardingRule]{
		Resource: "Network Load Balancer Forwarding Rule",
		List: func() ([]ionoscloud.NetworkLoadBalancerForwardingRule, error) {
			nlbForwardingRules, _, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(dcId, loadBalancerId)
			if err != nil {
				return nil, err
			}
			items, ok := nlbForwardingRules.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Network Load Balancer Forwarding Rules")
			}
			return *items, nil
		},
		Summary: func(rule ionoscloud.NetworkLoadBalancerForwardingRule) string {
			summary := ""
			if props, ok := rule.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil {
					summary += *name
				}
				if ip, ok := props.GetListenerIpOk(); ok && ip != nil && *ip != "" {
					summary += fmt.Sprintf(" (listenerIp: %s)", *ip)
				}
			}
			if id, ok := rule.GetIdOk(); ok && id != nil {
				summary += fmt.Sprintf(" (id: %s)", *id)
			}
			return summary
		},
		ID: func(rule ionoscloud.NetworkLoadBalancerForwardingRule) string {
			if id := rule.GetId(); id != nil {
				return *id
			}
			return ""
		},
		Delete: func(rule ionoscloud.NetworkLoadBalancerForwardingRule) error {
			resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, *rule.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
