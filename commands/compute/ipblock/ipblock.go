package ipblock

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "IpBlockId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "Size", JSONPath: "properties.size", Default: true},
	{Name: "Ips", JSONPath: "properties.ips", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func IpblockCmd() *core.Command {
	ipblockCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "ipblock",
			Aliases: []string{"ip", "ipb"},
			Short:   "Reserve and manage blocks of static public IPv4 addresses",
			Long: `An IpBlock is a reservation of one or more static, public IPv4 addresses, held in a single ` + "`--location`" + ` (region). Reserved IPs are yours to keep until you delete the block: unlike the dynamic (DHCP-assigned, ephemeral) addresses a NIC receives by default, a reserved IP does not change when a server is powered off ("Power Stop"/deallocated) or when its NIC is removed. This makes IpBlocks the right choice for anything that needs a stable address (DNS records, firewall allow-lists, published endpoints).

Once reserved, the individual IPs from a block are assigned to consumers within any Virtual Data Center in the SAME location:
  - Server NICs, as the primary or an additional IP (see ` + "`" + `compute nic` + "`" + ` --ips)
  - NAT Gateways, Network/Application Load Balancers
  - IP-failover groups (a shared IP that floats between NICs for HA)

Key properties:
  - Location is region-bound and fixed at reservation time; a block can only serve resources in its own location and cannot be moved.
  - Size (how many IPs are reserved) is set at creation and cannot be resized afterwards - reserve a new block instead.
  - Only the block's Name can be changed after creation (` + "`" + `ipblock update` + "`" + `).

Reserved IPs are billed for as long as they are held. To see which resource currently occupies each IP in a block, use ` + "`" + `ionosctl compute ipconsumer list --ipblock-id <id>` + "`" + `; an IP that is still in use by a consumer cannot be freed until that consumer releases it.`,
			TraverseChildren: true,
		},
	}
	ipblockCmd.AddColsFlag(allCols)

	ipblockCmd.AddCommand(IpBlockListCmd())
	ipblockCmd.AddCommand(IpBlockGetCmd())
	ipblockCmd.AddCommand(IpBlockCreateCmd())
	ipblockCmd.AddCommand(IpBlockUpdateCmd())
	ipblockCmd.AddCommand(IpBlockDeleteCmd())

	return core.WithConfigOverride(ipblockCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
