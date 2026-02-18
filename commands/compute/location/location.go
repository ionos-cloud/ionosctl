package location

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultLocationCols = []string{"LocationId", "Name", "CpuFamily"}
	allLocationCols     = []string{"LocationId", "Name", "Features", "ImageAliases", "CpuFamily"}
)

func LocationCmd() *core.Command {
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc"},
			Short:            "Location Operations",
			Long:             "The sub-command of `ionosctl location` allows you to see information about locations available to create objects.",
			TraverseChildren: true,
		},
	}
	globalFlags := locationCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLocationCols, tabheaders.ColsMessage(allLocationCols))
	_ = viper.BindPFlag(core.GetFlagName(locationCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = locationCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLocationCols, cobra.ShellCompDirectiveNoFileComp
	})

	locationCmd.AddCommand(LocationListCmd())
	locationCmd.AddCommand(LocationGetCmd())
	locationCmd.AddCommand(CpuCmd())

	return core.WithConfigOverride(locationCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
