package volume

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerVolumeDetachCmd() *core.Command {
	detachVolume := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "server",
		Resource:  "volume",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a Volume from a Server (the Volume is kept)",
		LongDesc: `Use this command to detach a Volume from a Server, or detach all attached Volumes at once with --all. Depending on the Volume's HotUnplug capability, detaching may trigger a reboot of the Server.

This does NOT delete the Volume or its data: the Volume stays in the Virtual Data Center and can be re-attached later or attached to a different Server. To permanently remove it, delete it separately with ` + "`ionosctl compute volume delete`" + `.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`,
		Example:    "ionosctl compute server volume detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID",
		PreCmdRun:  PreRunDcServerVolumeDetach,
		CmdRun:     RunServerVolumeDetach,
		InitClient: true,
	})
	detachVolume.AddColsFlag(allVolumeCols)
	detachVolume.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach every Volume currently attached to the Server (an alternative to giving a single --volume-id). Volumes are kept, not deleted")

	return detachVolume
}
