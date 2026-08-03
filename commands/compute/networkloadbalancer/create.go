package networkloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func NetworkLoadBalancerCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "networkloadbalancer",
		Resource:  "networkloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Network Load Balancer",
		LongDesc: `Use this command to create a Network Load Balancer (NLB) in a Virtual Data Center. An NLB is a layer-4 (TCP) load balancer that bridges a listener LAN (where clients connect) and a target LAN (the private network of the backend VMs it balances).

After creation the NLB has no forwarding rules yet - it will not forward any traffic until you add at least one rule with ` + "`nlb rule create`" + ` and attach targets with ` + "`nlb rule target add`" + `.

Networking flags:
  --listener-lan: LAN clients connect through (inbound). Defaults to 2.
  --target-lan:   private LAN of the balanced backend VMs (outbound). Defaults to 1.
  --ips:          the NLB's own IP addresses on the listener LAN. For a PUBLIC NLB these must be customer-reserved public IPs; for a PRIVATE NLB they are private IPs. These are also the IPs a forwarding rule listens on.
  --private-ips:  the NLB's private IPs (with subnet mask, e.g. 10.0.0.5/24) on the target LAN, used to reach the backends. If omitted, the system assigns one with a /24 mask.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id`,
		Example: `# Create an NLB with the default listener/target LANs
ionosctl compute networkloadbalancer create --datacenter-id DATACENTER_ID --name "web-nlb"

# Create a public NLB, pinning the LANs and giving it a reserved public listener IP
ionosctl compute networkloadbalancer create --datacenter-id DATACENTER_ID --name "prod-nlb" --listener-lan 2 --target-lan 1 --ips 203.0.113.10 --private-ips 10.0.0.5/24`,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunNetworkLoadBalancerCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Network Load Balancer", "Name of the Network Load Balancer")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "ID of the listener LAN (inbound) where clients connect to the NLB")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "ID of the private target LAN (outbound) hosting the backend VMs the NLB balances")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "IP addresses of the NLB on the listener LAN, also used as forwarding-rule listener IPs. Must be customer-reserved public IPs for a public NLB, or private IPs for a private NLB")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "Private IP addresses (with subnet mask, e.g. 10.0.0.5/24) the NLB uses on the target LAN to reach backends. If omitted, an IP with a /24 mask is generated")

	return cmd
}
