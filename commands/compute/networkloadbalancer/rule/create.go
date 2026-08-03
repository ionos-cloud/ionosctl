package rule

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerForwardingRuleCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "rule",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to create a forwarding rule on a Network Load Balancer. The rule listens on --listener-ip (which must be one of the NLB's own IPs) and --listener-port, then balances accepted TCP connections across the targets you later add with ` + "`nlb rule target add`" + `.

Pick a balancing algorithm with --algorithm:
  ROUND_ROBIN (default), LEAST_CONNECTION, RANDOM, SOURCE_IP (client-IP affinity).

The health-check flags tune resilience: --connect-timeout bounds how long a new connection to a target may take, --client-timeout / --target-timeout bound client- and target-side inactivity, and --retries sets how many times to retry a target after a connection failure before considering it down.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Listener Ip
* Listener Port`,
		Example: `# Create a rule that listens on port 80 (defaults to ROUND_ROBIN)
ionosctl compute networkloadbalancer rule create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --listener-ip 203.0.113.10 --listener-port 80

# Create a rule with client-IP affinity and custom health-check timeouts
ionosctl compute networkloadbalancer rule create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --name "https" --listener-ip 203.0.113.10 --listener-port 443 --algorithm SOURCE_IP --connect-timeout 5000 --retries 5`,
		PreCmdRun:  PreRunNetworkLoadBalancerForwardingRuleCreate,
		CmdRun:     RunNetworkLoadBalancerForwardingRuleCreate,
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
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Forwarding Rule", "The name for the Forwarding Rule")
	cmd.AddIpFlag(cloudapiv6.ArgListenerIp, "", nil, "Inbound IP the rule listens on. Must be one of the NLB's own IPs (--ips)", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgListenerPort, "", "", "Inbound TCP port the rule listens on. Range: 1 to 65535", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] Number of times to retry a target after a connection failure before marking it down. Range: 0 to 65535")
	cmd.AddIntFlag(cloudapiv6.ArgClientTimeout, "", 5000, "[Health Check] Maximum client-side inactivity, in milliseconds, before the connection is closed (client expected to acknowledge or send data)")
	cmd.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] Maximum time, in milliseconds, to wait for a connection to a target VM to succeed")
	cmd.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 5000, "[Health Check] Maximum target-side inactivity, in milliseconds, before the connection is closed")
	cmd.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Balancing algorithm used to pick a target per connection: ROUND_ROBIN, LEAST_CONNECTION, RANDOM, or SOURCE_IP (client-IP affinity)")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "RANDOM", "SOURCE_IP", "LEAST_CONNECTION"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
