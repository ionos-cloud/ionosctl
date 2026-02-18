package gpu

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ServerGpuListCmd() *core.Command {
	listGpus := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "gpu",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Gpus from a Server",
		LongDesc:   "List Gpus from a Server\n\nUse this command to retrieve a list of Gpus attached to a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    "ionosctl server gpu list --datacenter-id DATACENTER_ID --server-id SERVER_ID",
		PreCmdRun:  PreRunServerGpusList,
		CmdRun:     RunServerGpusList,
		InitClient: true,
	})
	listGpus.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listGpus.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	listGpus.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listGpus.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FilteredByTypeServersIds(viper.GetString(core.GetFlagName(listGpus.NS, cloudapiv6.ArgDataCenterId)), "GPU"), cobra.ShellCompDirectiveNoFileComp
	})

	return listGpus
}
