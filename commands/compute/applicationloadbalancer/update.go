package applicationloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ApplicationLoadBalancerUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an Application Load Balancer",
		LongDesc: `Use this command to update a specified Application Load Balancer from a Virtual Data Center. You can rename it, move which LANs it listens on / balances to, or change its public/private IP set. Only the flags you provide are changed.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Application Load Balancer Id`,
		Example: `# Rename an ALB
ionosctl compute applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --name "prod-web-alb"

# Replace the public listener IP set
ionosctl compute applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --ips 192.0.2.20 --wait`,
		PreCmdRun:  PreRunDcApplicationLoadBalancerIds,
		CmdRun:     RunApplicationLoadBalancerUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.ApplicationLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgApplicationLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ApplicationLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.Name(), cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Application Load Balancer", "The name of the Application Load Balancer.")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 0, "Numeric ID of the LAN clients connect to (the inbound/listener LAN).")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 0, "Numeric ID of the private LAN where the balanced backend servers live (the outbound/target LAN).")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "The IP addresses clients use to reach the balancer on the listener LAN. Customer-reserved public IPs for a public ALB, or private IPs for a private ALB. Replaces the existing set.")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "The balancer's own private IP addresses (with subnet mask) on the target LAN, used to reach the backends, e.g. --private-ips 10.0.1.5/24. If omitted, the system auto-generates an IP with a /24 subnet.")

	return cmd
}
