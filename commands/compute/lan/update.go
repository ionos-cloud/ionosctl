package lan

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LanUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "lan",
		Resource:  "lan",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a LAN",
		LongDesc: `Use this command to update an existing LAN. Only the flags you pass are changed; the rest are left as-is.

You can:
  * rename the LAN (` + "`--name`" + `);
  * switch it between public and private (` + "`--public`" + `). Turning a LAN public attaches it to an internet gateway (NICs can then get public IPv4s); turning it private removes that internet route;
  * attach it to a Cross-Connect (` + "`--pcc`" + `) to bridge it with private LANs in other VDCs of the same region. This requires the LAN to be private and to have a non-overlapping IP range with the other members. To detach from a Cross-Connect, set an empty value.

NOTE: IP failover groups are configured on a NIC, not here.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* LAN Id`,
		Example: `# Rename a LAN
ionosctl compute lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --name "renamed"

# Make a LAN private and attach it to a Cross-Connect to bridge it across VDCs
ionosctl compute lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --public=false --pcc PCC_ID`,
		PreCmdRun:  PreRunDcLanIds,
		CmdRun:     RunLanUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLanId, cloudapiv6.ArgIdShort, "", cloudapiv6.LanId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "New human-friendly name for the LAN")
	cmd.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", "ID of the Cross-Connect (Private Cross-Connect) to attach this LAN to, bridging it with private LANs in other VDCs of the same region. The LAN must be private and its IP range must not overlap the other members. Set an empty value to detach")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgPublic, "", cloudapiv6.DefaultPublic, "Whether the LAN is public. true = attached to an internet gateway so NICs can reach the internet and get public IPv4 addresses; false = private, internal traffic only. A LAN on a Cross-Connect (--pcc) must be private. E.g.: --public=true")

	cmd.AddStringFlag(cloudapiv6.FlagIPv6CidrBlock, "", "DISABLE", cloudapiv6.FlagIPv6CidrBlockDescriptionForLAN+
		` NOTE: Using an explicit Cidr to update the resource is not fully supported yet.`)

	return cmd
}
