package resource

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultResourceCols = []string{"ResourceId", "Name", "SecAuthProtection", "Type", "State"}
)

func ResourceCmd() *core.Command {
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Aliases:          []string{"res"},
			Short:            "Resource Operations",
			Long:             "The sub-commands of `ionosctl compute resource` allow you to list, get Resources.",
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultResourceCols, tabheaders.ColsMessage(defaultResourceCols))
	_ = viper.BindPFlag(core.GetFlagName(resourceCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})

	resourceCmd.AddCommand(ResourceListCmd())
	resourceCmd.AddCommand(ResourceGetCmd())

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
