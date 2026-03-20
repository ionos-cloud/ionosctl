package gpu

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

func ServerGpuGetCmd() *core.Command {
	getGpuCmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "gpu",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a GPU from a Server",
		LongDesc:   "Use this command to retrieve information about a GPU attached to a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* GPU Id",
		Example:    "ionosctl compute server gpu get --datacenter-id DATACENTER_ID --server-id SERVER_ID --gpu-id GPU_ID",
		PreCmdRun:  PreRunDcServerGpuIds,
		CmdRun:     RunServerGpuGet,
		InitClient: true,
	})
	getGpuCmd.AddStringSliceFlag(constants.ArgCols, "", defaultGpuCols, tabheaders.ColsMessage(allGpuCols))
	_ = getGpuCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allGpuCols, cobra.ShellCompDirectiveNoFileComp
	})
	getGpuCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = getGpuCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	getGpuCmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = getGpuCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FilteredByTypeServersIds(viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgDataCenterId)), "GPU"), cobra.ShellCompDirectiveNoFileComp
	})
	getGpuCmd.AddUUIDFlag(cloudapiv6.ArgGpuId, "", "", cloudapiv6.GpuId, core.RequiredFlagOption())
	_ = getGpuCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGpuId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GpusIds(viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	return getGpuCmd
}
