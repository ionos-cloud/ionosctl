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

func ServerVolumeListCmd() *core.Command {
	listVolumes := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "volume",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List attached Volumes from a Server",
		LongDesc:   "Use this command to retrieve a list of Volumes attached to the Server.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.VolumesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    "ionosctl server volume list --datacenter-id DATACENTER_ID --server-id SERVER_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunServerVolumeList,
		CmdRun:     cloudapiv6cmds.RunServerVolumesList,
		InitClient: true,
	})
	listVolumes.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = listVolumes.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(listVolumes.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	return listVolumes
}
