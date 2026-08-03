package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		LongDesc: `Use this command to update a forwarding rule's listener IP/port, balancing algorithm, or health-check settings. The rule's targets are managed separately via ` + "`nlb rule target`" + `.

Note that --listener-ip and --listener-port are required even when unchanged, and changing them moves the listener to a new IP/port pair (existing connections on the old pair are dropped).

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id`,
		Example: `# Switch a rule to the least-connection algorithm
ionosctl compute nlb rule update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID --listener-ip 203.0.113.10 --listener-port 80 --algorithm LEAST_CONNECTION

# Loosen health-check timeouts and retries
ionosctl compute nlb rule update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID --listener-ip 203.0.113.10 --listener-port 80 --connect-timeout 10000 --retries 5`,
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
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Inbound IP the rule listens on. Must be one of the NLB's own IPs (--ips)", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgListenerPort, "", "", "Inbound TCP port the rule listens on. Range: 1 to 65535", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Balancing algorithm used to pick a target per connection: ROUND_ROBIN, LEAST_CONNECTION, RANDOM, or SOURCE_IP (client-IP affinity)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] Number of times to retry a target after a connection failure before marking it down. Range: 0 to 65535")
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 5000, "[Health Check] Maximum client-side inactivity, in milliseconds, before the connection is closed (client expected to acknowledge or send data)")
	cmd.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] Maximum time, in milliseconds, to wait for a connection to a target VM to succeed")
	cmd.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 5000, "[Health Check] Maximum target-side inactivity, in milliseconds, before the connection is closed")

	return cmd
}
