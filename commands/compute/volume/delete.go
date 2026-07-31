package volume

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VolumeDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Volume",
		LongDesc: `Permanently delete a Volume from your Virtual Data Center. This destroys the block device and all data on it - the operation is irreversible, so use it with caution. If the Volume is currently attached to a Server, detach it first. DAS volumes bundled with Cube instances cannot be deleted independently.

Pass ` + "`" + `--all` + "`" + ` to delete every Volume in the datacenter. Use ` + "`" + `--force` + "`" + ` to skip the interactive confirmation prompt (useful in scripts), and ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the deletion has completed.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example:    `ionosctl compute volume delete --datacenter-id DATACENTER_ID --volume-id VOLUME_ID`,
		PreCmdRun:  PreRunDcVolumeDelete,
		CmdRun:     RunVolumeDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every Volume in the specified Datacenter. Combine with --force to skip per-volume confirmation")

	return cmd
}
