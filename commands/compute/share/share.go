package share

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type"}
	allGroupShareCols     = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type", "GroupId"}
)

func ShareCmd() *core.Command {
	shareCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "share",
			Short:            "Resource Share Operations",
			Long:             "The sub-commands of `ionosctl share` allow you to list, get, create, update, delete Resource Shares.",
			TraverseChildren: true,
		},
	}
	globalFlags := shareCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultGroupShareCols, tabheaders.ColsMessage(allGroupShareCols))
	_ = viper.BindPFlag(core.GetFlagName(shareCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = shareCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allGroupShareCols, cobra.ShellCompDirectiveNoFileComp
	})

	shareCmd.AddCommand(ShareListCmd())
	shareCmd.AddCommand(ShareGetCmd())
	shareCmd.AddCommand(ShareCreateCmd())
	shareCmd.AddCommand(ShareUpdateCmd())
	shareCmd.AddCommand(ShareDeleteCmd())

	return core.WithConfigOverride(shareCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
