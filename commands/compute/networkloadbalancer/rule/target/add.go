package target

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NlbRuleTargetAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Network Load Balancer Forwarding Rule Target",
		LongDesc: `Use this command to add a backend target to a forwarding rule. A target is a VM identified by --ip and --port on the NLB's target LAN; once added, the rule starts distributing connections to it according to the rule's balancing algorithm.

Weight (--weight): traffic is distributed in proportion to a target's weight relative to the sum of all targets' weights, so a higher weight means a higher share of connections. Default is 1, maximum is 256. A weight of 0 excludes the target from balancing but still lets it accept persistent connections. When sizing by capacity, start with mid-range values (e.g. 10-100) so you can adjust up or down later.

Health check (--check): when on (the default), the NLB periodically opens a TCP connection to the target's IP+port; the target is only considered available while it accepts these probes. When off, the target is always considered available. Use --check-interval to set the probe frequency, and --maintenance to force a target "down" (drain it) without removing it.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port`,
		Example: `# Add a backend VM with default weight and health checks
ionosctl compute networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip 10.0.0.11 --port 80

# Add a higher-capacity backend (double share) with a faster health-check probe
ionosctl compute networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip 10.0.0.12 --port 80 --weight 200 --check-interval 1000`,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleTarget,
		CmdRun:     RunNlbRuleTargetAdd,
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
	cmd.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of the backend target VM on the NLB's target LAN", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the backend target service. Range: 1 to 65535", core.RequiredFlagOption())
	cmd.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "Share of traffic this target receives relative to the other targets' weights. Range: 0 to 256; 0 excludes it from balancing but still accepts persistent connections")
	cmd.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] Interval in milliseconds between consecutive TCP health probes to the target")
	cmd.AddBoolFlag(cloudapiv6.ArgCheck, "", true, "[Health Check] When true, the target is only used while it accepts periodic TCP health probes; when false it is always considered available")
	cmd.AddBoolFlag(cloudapiv6.ArgMaintenance, "", false, "[Health Check] When true, drains the target: it is treated as down and receives no balanced traffic, even if healthy")

	return cmd
}
