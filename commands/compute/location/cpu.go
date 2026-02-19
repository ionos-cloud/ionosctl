package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultCpuCols = []string{"CpuFamily", "MaxCores", "MaxRam", "Vendor"}
)

func CpuCmd() *core.Command {
	cpuCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cpu",
			Short:            "Location CPU Architecture Operations",
			Long:             "The sub-command of `ionosctl compute location cpu` allows you to see information about available CPU Architectures in different Locations.",
			TraverseChildren: true,
		},
	}
	globalFlags := cpuCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultCpuCols, tabheaders.ColsMessage(defaultCpuCols))
	_ = viper.BindPFlag(core.GetFlagName(cpuCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = cpuCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultCpuCols, cobra.ShellCompDirectiveNoFileComp
	})

	cpuCmd.AddCommand(CpuListCmd())

	return core.WithConfigOverride(cpuCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func CpuListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "cpu",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List available CPU Architecture from a Location",
		LongDesc:   "Use this command to get information about available CPU Architectures from a specific Location.\n\nRequired values to run command:\n\n* Location Id",
		Example:    `ionosctl compute location cpu list --location-id LOCATION_ID`,
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationCpuList,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLocationId, "", "", cloudapiv6.LocationId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
