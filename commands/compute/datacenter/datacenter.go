package datacenter

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultDatacenterCols = []string{"DatacenterId", "Name", "Location", "CpuFamily", "IPv6CidrBlock", "State"}
	allDatacenterCols     = []string{"DatacenterId", "Name", "Location", "State", "Description", "Version",
		"Features", "CpuFamily", "SecAuthProtection", "IPv6CidrBlock"}
)

func DatacenterCmd() *core.Command {
	datacenterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Aliases:          []string{"d", "dc", "vdc"},
			Args:             cobra.ExactValidArgs(1),
			Short:            "Data Center Operations",
			Long:             "The sub-commands of `ionosctl datacenter` allow you to create, list, get, update and delete Data Centers.",
			TraverseChildren: true,
		},
	}
	globalFlags := datacenterCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultDatacenterCols, tabheaders.ColsMessage(allDatacenterCols))
	_ = viper.BindPFlag(core.GetFlagName(datacenterCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = datacenterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allDatacenterCols, cobra.ShellCompDirectiveNoFileComp
	})

	datacenterCmd.AddCommand(DatacenterListCmd())
	datacenterCmd.AddCommand(DatacenterGetCmd())
	datacenterCmd.AddCommand(DatacenterCreateCmd())
	datacenterCmd.AddCommand(DatacenterUpdateCmd())
	datacenterCmd.AddCommand(DatacenterDeleteCmd())

	return core.WithConfigOverride(datacenterCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
