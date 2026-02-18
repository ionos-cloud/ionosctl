package volume

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
		ShortDesc: "Detach a Volume from a Server",
		LongDesc: `This will detach the Volume from the Server. Depending on the Volume HotUnplug settings, this may result in the Server being rebooted. This will NOT delete the Volume from your Virtual Data Center. You will need to use a separate command to delete a Volume.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`,
		Example:    "ionosctl server volume detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunDcServerVolumeDetach,
		CmdRun:     cloudapiv6cmds.RunServerVolumeDetach,
		InitClient: true,
	})
	detachVolume.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = detachVolume.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
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
	detachVolume.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume detachment to be executed")
	detachVolume.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume detachment [seconds]")
	detachVolume.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach all Volumes.")

	return detachVolume
}
