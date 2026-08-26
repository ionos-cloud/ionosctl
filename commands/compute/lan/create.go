package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func LanCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a LAN",
		LongDesc: `Use this command to create a new LAN inside a Virtual Data Center (VDC). A LAN is a virtual network segment that connects the NICs attached to it within a single datacenter (` + "`--datacenter-id`" + `).

Decide at creation time whether the LAN is:
  * private (` + "`--public=false`" + `, the default): internal, datacenter-local traffic only, no direct internet route. This is the type required to later join a Cross-Connect.
  * public (` + "`--public=true`" + `): attached to an internet gateway; NICs on it can send/receive public internet traffic and be assigned public IPv4 addresses.

Optionally pass ` + "`--pcc-id`" + ` to immediately attach the LAN to a Cross-Connect (Private Cross-Connect), bridging it with private LANs in other VDCs in the same region. The LAN must be private and its IP range must not overlap the other members of the Cross-Connect.

NOTE: IP failover groups (a reserved IP that floats between servers for high availability) are NOT set here; configure them on a NIC after the LAN exists.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id`,
		Example: `# Create a simple private LAN (internal traffic only)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "backend"

# Create a public LAN (internet-facing, NICs can receive public IPs)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "frontend" --public=true

# Create a private LAN and attach it to a Cross-Connect to bridge it with LANs in other VDCs (Cross-Connect requires a private LAN)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "cross-dc" --public=false --pcc-id PCC_ID`,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunLanCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed LAN", "A human-friendly name for the LAN. Not required to be unique")
	cmd.AddBoolFlag(cloudapiv6.ArgPublic, cloudapiv6.ArgPublicShort, cloudapiv6.DefaultPublic, "Whether the LAN is public. true = attached to an internet gateway so NICs can reach the internet and be assigned public IPv4 addresses; false (default) = private, internal datacenter traffic only. A Cross-Connect (--pcc-id) requires a private LAN. E.g.: --public=true")
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", "ID of the Cross-Connect (Private Cross-Connect) to attach this LAN to, bridging it with private LANs in other VDCs of the same region. The LAN must be private (--public=false) and its IP range must not overlap the other Cross-Connect members")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(cloudapiv6.FlagIPv6CidrBlock, "", "DISABLE", cloudapiv6.FlagIPv6CidrBlockDescriptionForLAN+
		` Use "DISABLE" (default) to keep the LAN IPv4-only, "AUTO" to let IONOS assign a /64 block automatically from the datacenter's range, or pass an explicit /64 block within the datacenter's IPv6 range.`)

	return cmd
}
