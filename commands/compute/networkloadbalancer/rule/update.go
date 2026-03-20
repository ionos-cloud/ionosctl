package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerForwardingRuleUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to update a specified Network Load Balancer Forwarding Rule from a Network Load Balancer. You can also update Health Check settings.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example:    `ionosctl compute nlb rule update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME`,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgRuleId, cloudapiv6.ArgIdShort, "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "The name for the Forwarding Rule")
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Listening IP", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgListenerPort, "", "", "Listening port number. Range: 1 to 65535", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] Retries specifies the number of retries to perform on a target VM after a connection failure. Range: 0 to 65535")
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 5000, "[Health Check] ClientTimeout is expressed in milliseconds. This inactivity timeout applies when the client is expected to acknowledge or send data")
	cmd.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	cmd.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 5000, "[Health Check] TargetTimeout specifies the maximum inactivity time (in milliseconds) on the target VM side")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Forwarding Rule update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule update [seconds]")

	return cmd
}
