package rule

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
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
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
	c.Verbose("Network Load Balancer Forwarding Rule with id: %v is getting...", c.Flags().String(cloudapiv6.ArgRuleId))

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgRuleId),
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
		proper.SetAlgorithm(c.Flags().String(cloudapiv6.ArgAlgorithm))
	}

	if !proper.HasName() {
		proper.SetName(c.Flags().String(cloudapiv6.ArgName))
	}

	health := getForwardingRuleHealthCheckPropertiesSet(c)
	if health != nil {
		proper.SetHealthCheck(health.NetworkLoadBalancerForwardingRuleHealthCheck)
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateForwardingRule(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
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
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgRuleId),
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
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	loadBalancerId := c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId)
	ruleId := c.Flags().String(cloudapiv6.ArgRuleId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
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

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgAlgorithm) {
		algorithm := strings.ToUpper(c.Flags().String(cloudapiv6.ArgAlgorithm))
		input.SetAlgorithm(algorithm)

		c.Verbose("Property Algorithm set: %v", algorithm)
	}

	if c.Flags().Changed(cloudapiv6.ArgListenerIp) {
		listenerIp := c.Flags().String(cloudapiv6.ArgListenerIp)
		input.SetListenerIp(listenerIp)

		c.Verbose("Property ListenerIp set: %v", listenerIp)
	}

	if c.Flags().Changed(cloudapiv6.ArgListenerPort) {
		listenerPort := c.Flags().Int32(cloudapiv6.ArgListenerPort)
		input.SetListenerPort(listenerPort)

		c.Verbose("Property ListenerPort set: %v", listenerPort)
	}

	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: input,
	}
}

func getForwardingRuleHealthCheckPropertiesSet(c *core.CommandConfig) *resources.NetworkLoadBalancerForwardingRuleHealthCheck {
	inputHealth := ionoscloud.NetworkLoadBalancerForwardingRuleHealthCheck{}

	if c.Flags().Changed(cloudapiv6.ArgRetries) {
		inputHealth.SetRetries(c.Flags().Int32(cloudapiv6.ArgRetries))
	}

	if c.Flags().Changed(cloudapiv6.ArgClientTimeout) {
		inputHealth.SetClientTimeout(c.Flags().Int32(cloudapiv6.ArgClientTimeout))
	}

	if c.Flags().Changed(cloudapiv6.ArgConnectionTimeout) {
		inputHealth.SetConnectTimeout(c.Flags().Int32(cloudapiv6.ArgConnectionTimeout))
	}

	if c.Flags().Changed(cloudapiv6.ArgTargetTimeout) {
		inputHealth.SetTargetTimeout(c.Flags().Int32(cloudapiv6.ArgTargetTimeout))
	}

	return &resources.NetworkLoadBalancerForwardingRuleHealthCheck{
		NetworkLoadBalancerForwardingRuleHealthCheck: inputHealth,
	}
}

func DeleteAllNetworkLoadBalancerForwardingRules(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	loadBalancerId := c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Network Load Balancer ID: %v", loadBalancerId)
	c.Verbose("Getting Network Load Balancer Forwarding Rules...")

	nlbForwardingRules, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListForwardingRules(dcId, loadBalancerId)
	if err != nil {
		return err
	}

	nlbForwardingRulesItems, ok := nlbForwardingRules.GetItemsOk()
	if !ok || nlbForwardingRulesItems == nil {
		return fmt.Errorf("could not get items of Network Load Balancer Forwarding Rules")
	}

	if len(*nlbForwardingRulesItems) <= 0 {
		return fmt.Errorf("no Network Load Balancer Forwarding Rules found")
	}

	var multiErr error
	for _, nlbForwardingRule := range *nlbForwardingRulesItems {
		name := nlbForwardingRule.GetProperties().Name
		id := nlbForwardingRule.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Network Load Balancer Forwarding Rule with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NetworkLoadBalancers().DeleteForwardingRule(dcId, loadBalancerId, *id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
