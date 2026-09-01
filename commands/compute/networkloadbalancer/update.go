package networkloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NetworkLoadBalancerUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Network Load Balancer",
		LongDesc: `Use this command to update a Network Load Balancer's name, listener/target LANs, or IP addresses. Only the flags you pass are changed.

Note that changing --listener-lan, --target-lan, or the IP flags rewires the NLB's connectivity and can interrupt in-flight traffic; forwarding rules that listen on an IP you remove will stop serving.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Network Load Balancer Id`,
		Example: `# Rename an NLB
ionosctl compute networkloadbalancer update --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID --name "renamed-nlb"

# Repoint the NLB to a different target LAN
ionosctl compute networkloadbalancer update --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID --target-lan 3`,
		PreCmdRun:  PreRunDcNetworkLoadBalancerIds,
		CmdRun:     RunNetworkLoadBalancerUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgIdShort, "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "ID of the listener LAN (inbound) where clients connect to the NLB")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "ID of the private target LAN (outbound) hosting the backend VMs the NLB balances")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "IP addresses of the NLB on the listener LAN, also used as forwarding-rule listener IPs. Must be customer-reserved public IPs for a public NLB, or private IPs for a private NLB")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Private IP addresses (with subnet mask, e.g. 10.0.0.5/24) the NLB uses on the target LAN to reach backends. If omitted, an IP with a /24 mask is generated")

	return cmd
}
