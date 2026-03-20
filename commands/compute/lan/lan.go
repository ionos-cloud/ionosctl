package lan

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultLanCols = []string{"LanId", "Name", "Public", "PccId", "IPv6CidrBlock", "State"}
	allLanCols     = []string{"LanId", "Name", "Public", "PccId", "IPv6CidrBlock", "State", "DatacenterId"}
)

func LanCmd() *core.Command {
	lanCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "lan",
			Aliases:          []string{"l"},
			Short:            "LAN Operations",
			Long:             "The sub-commands of `ionosctl compute lan` allow you to create, list, get, update, delete LANs.",
			TraverseChildren: true,
		},
	}
	globalFlags := lanCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLanCols, tabheaders.ColsMessage(allLanCols))
	_ = viper.BindPFlag(core.GetFlagName(lanCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = lanCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLanCols, cobra.ShellCompDirectiveNoFileComp
	})

	lanCmd.AddCommand(LanListCmd())
	lanCmd.AddCommand(LanGetCmd())
	lanCmd.AddCommand(LanCreateCmd())
	lanCmd.AddCommand(LanUpdateCmd())
	lanCmd.AddCommand(LanDeleteCmd())

	return core.WithConfigOverride(lanCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
