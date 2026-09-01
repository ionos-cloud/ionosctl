package ipblock

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func IpBlockCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Reserve a block of static public IPv4 addresses",
		LongDesc: `Reserve an IpBlock: a set of ` + "`--size`" + ` static, public IPv4 addresses held in one ` + "`--location`" + ` (region). The reserved IPs can then be assigned to NICs, NAT gateways, load balancers or IP-failover groups in any Virtual Data Center in that SAME location.

Both --location and --size are fixed at reservation time and cannot be changed later - a block cannot be moved to another region or resized. Reserved IPs are billed while held and persist across server power-offs (unlike dynamic DHCP addresses), so reserve only what you need.

Reservation is asynchronous. Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the IpBlock reaches AVAILABLE state and its IPs are ready to assign.`,
		Example: `# Reserve a block of 1 IP in Berlin (de/txl)
ionosctl compute ipblock create --name web-vip --location de/txl --size 1 --wait

# Reserve 4 IPs in Frankfurt for a load balancer / failover pool, then list the assigned addresses
ionosctl compute ipblock create --name lb-pool --location de/fra --size 4 --wait`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunIpBlockCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "A friendly label for the block (shown in listings; not the IP addresses themselves). If omitted, the API assigns a name automatically")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "de/txl", "Region the IPs are reserved in, as <country>/<city> (e.g. de/txl Berlin, de/fra Frankfurt, us/las Las Vegas). The block can only serve resources in this region and cannot be moved afterwards. Location de/fra/2 is currently unavailable")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgSize, "", 2, "How many static public IPv4 addresses to reserve in this block. Fixed at creation - a block cannot be resized later. Each reserved IP is billed while held")

	return cmd
}
