package applicationloadbalancer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ApplicationLoadBalancerCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "applicationloadbalancer",
		Resource:  "applicationloadbalancer",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an Application Load Balancer",
		LongDesc: `Use this command to create a layer-7 Application Load Balancer in a specified Virtual Data Center.

An ALB bridges two LANs in the same data center:
  * The listener LAN (--listener-lan) is where clients reach the balancer. For a public ALB this LAN is internet-facing and --ips holds customer-reserved public IPs; for a private ALB it is an internal LAN and --ips holds private IPs.
  * The target LAN (--target-lan) is the private LAN where your backend servers (grouped into target groups) live. The ALB uses --private-ips as its own addresses on this LAN to reach the backends.

After creating the ALB you attach forwarding rules (` + "`" + `alb rule create` + "`" + `), then HTTP rules within them (` + "`" + `alb rule httprule add` + "`" + `) to route traffic to target groups.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the ALB reaches AVAILABLE state before its rules are configured.

Required values to run command:

* Data Center Id`,
		Example: `# Create a public ALB listening on LAN 2, balancing to backends on LAN 1
ionosctl compute applicationloadbalancer create --datacenter-id DATACENTER_ID --name "web-alb" --listener-lan 2 --target-lan 1 --ips 192.0.2.10

# Create an ALB and wait for it to become AVAILABLE, letting the system auto-assign a /24 private IP on the target LAN
ionosctl compute applicationloadbalancer create --datacenter-id DATACENTER_ID --name "web-alb" --listener-lan 2 --target-lan 1 --ips 192.0.2.10 --private-ips 10.0.1.5/24 --wait`,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunApplicationLoadBalancerCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Application Load Balancer", "The name of the Application Load Balancer.")
	cmd.AddIntFlag(cloudapiv6.ArgListenerLan, "", 2, "Numeric ID of the LAN clients connect to (the inbound/listener LAN). For a public ALB this is an internet-facing LAN; for a private ALB it is an internal LAN. Defaults to 2.")
	cmd.AddIntFlag(cloudapiv6.ArgTargetLan, "", 1, "Numeric ID of the private LAN where the balanced backend servers live (the outbound/target LAN). Defaults to 1.")
	cmd.AddStringSliceFlag(cloudapiv6.ArgIps, "", nil, "The IP addresses clients use to reach the balancer on the listener LAN. These are customer-reserved public IPs for a public ALB, or private IPs for a private ALB. Provide one or more, e.g. --ips 192.0.2.10,192.0.2.11")
	cmd.AddStringSliceFlag(cloudapiv6.ArgPrivateIps, "", nil, "The balancer's own private IP addresses (with subnet mask) on the target LAN, used to reach the backends. Each value must include a valid subnet mask, e.g. --private-ips 10.0.1.5/24. If omitted, the system auto-generates an IP with a /24 subnet.")

	return cmd
}
