package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Snapshot of a Volume within the Virtual Data Center",
		LongDesc: `Use this command to take a Snapshot of a storage Volume. A Snapshot is created from the perspective of a specific Volume, so both the Volume Id and the Data Center Id that Volume lives in are required.

The Snapshot captures the FULL provisioned capacity of the Volume (not just the used space) as a point-in-time image, and is stored independently at the Volume's LOCATION. From then on it can be restored onto a Volume (` + "`snapshot restore`" + `) or used as a boot image when creating new Volumes - but only within that same location and your contract.

For a consistent image, prefer taking the Snapshot while the source Volume's Server is powered off (or with in-guest I/O quiesced); a Snapshot of a live, busy filesystem may be crash-consistent only.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the Snapshot to reach AVAILABLE state before the command returns.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example: `# Take a basic snapshot of a volume
ionosctl compute snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name "pre-upgrade baseline"

# Advanced: name it, note the OS licence, and require Contract Owner / re-authentication to delete or restore it
ionosctl compute snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name "prod-db golden" --description "before schema migration" --licence-type LINUX --sec-auth-protection=true --wait`,
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunSnapshotCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Snapshot", "A human-friendly label for the Snapshot; shown in listings. Does not have to be unique")
	cmd.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Free-form notes about the Snapshot, e.g. why or when it was taken")
	cmd.AddSetFlag(cloudapiv6.ArgLicenceType, "", "LINUX", constants.EnumLicenceType, "The operating-system licence recorded on the Snapshot. This carries over to Volumes created/restored from it and affects OS licensing (notably WINDOWS variants are billed). Set it to match the source Volume's OS")
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgSecAuthProtection, "", false, "Protect the Snapshot with secure authentication: when true, deleting or restoring it requires the Contract Owner or a re-authenticated user, guarding against accidental loss. E.g.: --sec-auth-protection=true, --sec-auth-protection=false")

	return cmd
}
