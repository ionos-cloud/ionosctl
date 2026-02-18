package target

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NlbRuleTargetRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
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
		Example:    `ionosctl nlb rule target remove --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip TARGET_IP --port TARGET_PORT`,
		PreCmdRun:  cloudapiv6cmds.PreRunNetworkLoadBalancerRuleTargetRemove,
		CmdRun:     cloudapiv6cmds.RunNlbRuleTargetRemove,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of a balanced target VM", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the balanced target service. Range: 1 to 65535", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule Target deletion to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Target deletion [seconds]")
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Forwarding Rule Targets.")

	return cmd
}
