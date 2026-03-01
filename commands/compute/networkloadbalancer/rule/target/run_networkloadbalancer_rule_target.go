package target

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunNetworkLoadBalancerRuleTarget(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgIp, cloudapiv6.ArgPort)
}

func PreRunNetworkLoadBalancerRuleTargetRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgTargetIp, cloudapiv6.ArgTargetPort},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgAll},
	)
}

func RunNlbRuleTargetList(c *core.CommandConfig) error {
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := ng.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting rule properties")
	}

	targets, ok := properties.GetTargetsOk()
	if !ok || targets == nil {
		return fmt.Errorf("error getting rule targets")
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NetworkLoadBalancerRuleTarget, *targets,
		tabheaders.GetHeadersAllDefault(defaultRuleTargetCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", out)

	return nil
}

func RunNlbRuleTargetAdd(c *core.CommandConfig) error {
	var targetItems []ionoscloud.NetworkLoadBalancerForwardingRuleTarget

	ngOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if properties, ok := ngOld.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			targetItems = *targets
		}
	}

	targetNew := getRuleTargetInfo(c)
	targetItems = append(targetItems, targetNew.NetworkLoadBalancerForwardingRuleTarget)
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))

	nlbForwardingRule := &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Adding NlbRuleTarget with id: %v to NetworkLoadBalancer with id: %v", ruleId, nlbId))

	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbForwardingRule)
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.NetworkLoadBalancerRuleTarget,
		targetNew.NetworkLoadBalancerForwardingRuleTarget, tabheaders.GetHeadersAllDefault(defaultRuleTargetCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", out)

	return nil
}

func RunNlbRuleTargetRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllNlbRuleTarget(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete forwarding rule target", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"NlbRuleTarget with id: %v is removing...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))))

	frOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	proper, err := getRuleTargetsRemove(c, frOld)
	if err != nil {
		return fmt.Errorf("could not update Forwarding Rule properties for delete: %w", err)
	}

	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		proper,
	)
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Network Load Balancer Forwarding Rule Target successfully deleted"))
	return nil
}

func RemoveAllNlbRuleTarget(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("NetworkLoadBalancer ID: %v", nlbId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("NetworkLoadBalancerForwardingRule ID: %v", ruleId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting NetworkLoadBalancerForwardingRule..."))

	forwardingRule, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(dcId, nlbId, ruleId)
	if err != nil {
		return err
	}

	forwardingRuleProperties, ok := forwardingRule.GetPropertiesOk()
	if !ok || forwardingRuleProperties == nil {
		return fmt.Errorf("could not get Forwarding Rule properties")
	}

	targets, ok := forwardingRuleProperties.GetTargetsOk()
	if !ok || targets == nil {
		return fmt.Errorf("could not get items of Forwarding Rule Targets")
	}

	if len(*targets) <= 0 {
		return fmt.Errorf("no Forwarding Rule Targets found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Forwarding Rule Targets to be removed:"))

	for _, target := range *targets {
		delIdAndName := ""

		if ipOk, ok := target.GetIpOk(); ok && ipOk != nil {
			delIdAndName += " Forwarding Rule Target IP: " + *ipOk
		}

		if portOk, ok := target.GetPortOk(); ok && portOk != nil {
			delIdAndName += " Forwarding Rule Target Port: " + strconv.Itoa(int(*portOk))
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("%s", delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the Forwarding Rule Targets", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing all the Forwarding Rule Targets..."))

	targetItems := make([]ionoscloud.NetworkLoadBalancerForwardingRuleTarget, 0)
	nlbFwRuleProp := &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}

	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbFwRuleProp)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Network Load Balancer Forwarding Rule Targets successfully deleted"))
	return nil
}

func getRuleTargetInfo(c *core.CommandConfig) resources.NetworkLoadBalancerForwardingRuleTarget {
	targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	targetPort := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort))
	weight := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgWeight))
	target := resources.NetworkLoadBalancerForwardingRuleTarget{}

	target.SetIp(targetIp)
	target.SetPort(targetPort)
	target.SetWeight(weight)

	targetHealth := resources.NetworkLoadBalancerForwardingRuleTargetHealthCheck{}
	maintenance := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenance))
	check := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCheck))
	checkInterval := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))

	targetHealth.SetMaintenance(maintenance)
	targetHealth.SetCheck(check)
	targetHealth.SetCheckInterval(checkInterval)
	target.SetHealthCheck(targetHealth.NetworkLoadBalancerForwardingRuleTargetHealthCheck)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Properties set for adding the NlbRuleTarget: Ip: %v, Port: %v, Weight: %v, Maintenance: %v, Check: %v, CheckInterval: %v",
		targetIp, targetPort, weight, maintenance, check, checkInterval))

	return target
}

func getRuleTargetsRemove(c *core.CommandConfig, frOld *resources.NetworkLoadBalancerForwardingRule) (*resources.NetworkLoadBalancerForwardingRuleProperties, error) {
	var (
		foundIp   = false
		foundPort = false
	)

	targetItems := make([]ionoscloud.NetworkLoadBalancerForwardingRuleTarget, 0)

	properties, ok := frOld.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve old Forwarding Rule properties")
	}

	targets, ok := properties.GetTargetsOk()
	if !ok || targets == nil {
		return nil, fmt.Errorf("could not retrieve old Forwarding Rule Targets")
	}

	// Iterate through all targets
	for _, targetItem := range *targets {
		removeIp := false
		removePort := false

		if ip, ok := targetItem.GetIpOk(); ok && ip != nil {
			if *ip == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)) {
				removeIp = true
				foundIp = true
			}
		}

		if port, ok := targetItem.GetPortOk(); ok && port != nil {
			if *port == viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort)) {
				removePort = true
				foundPort = true
			}
		}

		if removeIp && removePort {
			continue
		}

		targetItems = append(targetItems, targetItem)
	}

	if !foundIp {
		return nil, errors.New("no forwarding rule target with the specified IP found")
	}

	if !foundPort {
		return nil, errors.New("no forwarding rule target with the specified port found")
	}

	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}, nil
}

func PreRunDcNetworkLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId)
}
