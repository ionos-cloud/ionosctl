package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotRestoreCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "snapshot",
		Resource:  "snapshot",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a Snapshot onto a Volume (destructive overwrite)",
		LongDesc: `Use this command to write the contents of a Snapshot back onto an existing target Volume, rolling that Volume to the point in time the Snapshot captured.

This is a DESTRUCTIVE, in-place overwrite: the target Volume's current data is replaced by the Snapshot image, so anything written since is lost. You are asked to confirm (use --force to skip the prompt). The Snapshot and the target Volume must be at the same LOCATION, and the target Volume should be at least as large as the Snapshot.

The target Volume is identified by its Data Center Id + Volume Id; it does NOT have to be the Volume the Snapshot was originally taken from - any compatible Volume at the same location works. To spin up a brand-new Volume from a Snapshot instead of overwriting one, pass the Snapshot Id as the image when creating a Volume rather than using restore.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait until the restore completes and the Volume is AVAILABLE.

Required values to run command:

* Datacenter Id
* Volume Id
* Snapshot Id`,
		Example: `# Roll a volume back to a snapshot (prompts for confirmation)
ionosctl compute snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --wait

# Advanced: restore onto a different target volume, no prompt, in a script
ionosctl compute snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id OTHER_VOLUME_ID --force --wait`,
		PreCmdRun:  PreRunSnapshotIdDcIdVolumeId,
		CmdRun:     RunSnapshotRestore,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgSnapshotId, cloudapiv6.ArgIdShort, "", cloudapiv6.SnapshotId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSnapshotId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.SnapshotIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, "", "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
