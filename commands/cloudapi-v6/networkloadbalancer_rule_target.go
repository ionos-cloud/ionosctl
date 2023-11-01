package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultRuleTargetCols = []string{"TargetIp", "TargetPort", "Weight", "Check", "CheckInterval", "Maintenance"}
)

func NlbRuleTargetCmd() *core.Command {
	ctx := context.TODO()
	nlbRuleTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Network Load Balancer Forwarding Rule Target Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer rule target` allow you to add, list, update, remove Network Load Balancer Forwarding Rule Targets.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleTargetCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultRuleTargetCols, tabheaders.ColsMessage(defaultRuleTargetCols))
	_ = viper.BindPFlag(core.GetFlagName(nlbRuleTargetCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = nlbRuleTargetCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultRuleTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace:  "forwardingrule",
		Resource:   "target",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancer Forwarding Rule Targets",
		LongDesc:   "Use this command to list Targets of a Network Load Balancer Forwarding Rule.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id\n* Forwarding Rule Id",
		Example:    listNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNlbRuleTargetList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Network Load Balancer Forwarding Rule Target",
		LongDesc: `Use this command to add a Forwarding Rule Target in a specified Network Load Balancer Forwarding Rule. You can also set Health Check Settings for Forwarding Rule Target. The Check parameter for Health Check Settings specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

Regarding the Weight parameter, this parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port`,
		Example:    addNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleTarget,
		CmdRun:     RunNlbRuleTargetAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of a balanced target VM", core.RequiredFlagOption())
	add.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the balanced target service. Range: 1 to 65535", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "Weight parameter is used to adjust the target VM's weight relative to other target VMs. Maximum: 256")
	add.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] CheckInterval determines the duration (in milliseconds) between consecutive health checks")
	add.AddBoolFlag(cloudapiv6.ArgCheck, "", true, "[Health Check] Check specifies whether the target VM's health is checked")
	add.AddBoolFlag(cloudapiv6.ArgMaintenance, "", false, "[Health Check]  Maintenance specifies if a target VM should be marked as down, even if it is not")
	add.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Target creation to be executed")
	add.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Target creation [seconds]")
	add.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "target",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Target from a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to remove a specified Target from Network Load Balancer Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port`,
		Example:    removeNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleTargetRemove,
		CmdRun:     RunNlbRuleTargetRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of a balanced target VM", core.RequiredFlagOption())
	removeCmd.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the balanced target service. Range: 1 to 65535", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Target deletion to be executed")
	removeCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Target deletion [seconds]")
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Forwarding Rule Targets.")
	removeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return nlbRuleTargetCmd
}

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
		resources.QueryParams{},
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), out)

	return nil
}

func RunNlbRuleTargetAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	var targetItems []ionoscloud.NetworkLoadBalancerForwardingRuleTarget

	ngOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding NlbRuleTarget with id: %v to NetworkLoadBalancer with id: %v", ruleId, nlbId))

	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbForwardingRule, queryParams)
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), out)

	return nil
}

func RunNlbRuleTargetRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllNlbRuleTarget(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete forwarding rule target", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"NlbRuleTarget with id: %v is removing...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))))

	frOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
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
		queryParams,
	)
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Network Load Balancer Forwarding Rule Target successfully deleted"))
	return nil
}

func RemoveAllNlbRuleTarget(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("NetworkLoadBalancer ID: %v", nlbId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("NetworkLoadBalancerForwardingRule ID: %v", ruleId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting NetworkLoadBalancerForwardingRule..."))

	forwardingRule, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(dcId, nlbId, ruleId, cloudapiv6.ParentResourceQueryParams)
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Forwarding Rule Targets to be removed:"))

	for _, target := range *targets {
		delIdAndName := ""

		if ipOk, ok := target.GetIpOk(); ok && ipOk != nil {
			delIdAndName += " Forwarding Rule Target IP: " + *ipOk
		}

		if portOk, ok := target.GetPortOk(); ok && portOk != nil {
			delIdAndName += " Forwarding Rule Target Port: " + strconv.Itoa(int(*portOk))
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the Forwarding Rule Targets", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing all the Forwarding Rule Targets..."))

	targetItems := make([]ionoscloud.NetworkLoadBalancerForwardingRuleTarget, 0)
	nlbFwRuleProp := &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}

	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbFwRuleProp, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Network Load Balancer Forwarding Rule Targets successfully deleted"))
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
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
