package volume

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerVolumeGetCmd() *core.Command {
	getVolumeCmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "volume",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an attached Volume from a Server",
		LongDesc:   "Use this command to retrieve information about an attached Volume on Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Volume Id",
		Example:    "ionosctl compute server volume get --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID",
		InitClient: true,
		PreCmdRun:  PreRunDcServerVolumeIds,
		CmdRun:     RunServerVolumeGet,
	})
	getVolumeCmd.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	return getVolumeCmd
}
