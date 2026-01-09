package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultGpuCols = []string{"GpuId", "Type", "Vendor", "Model", "Name", "State"}
	allGpuCols     = []string{"GpuId", "Type", "Vendor", "Model", "Name", "State"}
)

func ServerGpuCmd() *core.Command {
	ctx := context.TODO()
	serverGpuCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gpu",
			Aliases:          []string{"gpus"},
			Short:            "GPU operations",
			Long:             "The sub-commands of `ionosctl server gpu` allow you to get and list Gpus from a Server.",
			TraverseChildren: true,
		},
	}

	/*
		List GPUs command
	*/

	listGpus := core.NewCommand(ctx, serverGpuCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "gpu",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Gpus from a Server",
		LongDesc:   "List Gpus from a Server\n\nUse this command to retrieve a list of Gpus attached to a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listGpusServerExample,
		PreCmdRun:  PreRunServerGpusList,
		CmdRun:     RunServerGpusList,
		InitClient: true,
	})

	listGpus.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listGpus.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	listGpus.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listGpus.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FilteredByTypeServersIds(viper.GetString(core.GetFlagName(listGpus.NS, cloudapiv6.ArgDataCenterId)), serverGPUType), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get GPU Command
	*/
	getGpuCmd := core.NewCommand(ctx, serverGpuCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "gpu",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a GPU from a Server",
		LongDesc:   "Use this command to retrieve information about a GPU attached to a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* GPU Id",
		Example:    getGpuServerExample,
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
		return completer.FilteredByTypeServersIds(viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgDataCenterId)), serverGPUType), cobra.ShellCompDirectiveNoFileComp
	})
	getGpuCmd.AddUUIDFlag(cloudapiv6.ArgGpuId, "", "", cloudapiv6.GpuId, core.RequiredFlagOption())
	_ = getGpuCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGpuId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GpusIds(viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(getGpuCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(serverGpuCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func PreRunServerGpusList(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)

}

func RunServerGpusList(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	gpus, _, err := client.Must().CloudClient.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(c.Context, dcId, serverId).Execute()

	if err != nil {
		return fmt.Errorf("failed to list Gpus from Server: %w", err)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Gpu, gpus,
		tabheaders.GetHeaders(allGpuCols, defaultGpuCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func PreRunDcServerGpuIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgGpuId)
}

func RunServerGpuGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	gpuId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGpuId))

	gpu, _, err := client.Must().CloudClient.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(c.Context, dcId, serverId, gpuId).Execute()

	if err != nil {
		return fmt.Errorf("failed to get GPU from Server: %w", err)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Gpu, gpu,
		tabheaders.GetHeaders(allGpuCols, defaultGpuCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
