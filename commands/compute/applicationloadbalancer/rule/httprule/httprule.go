package httprule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultAlbRuleHttpRuleCols = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Condition"}
	allAlbRuleHttpRuleCols     = []string{"Name", "Type", "TargetGroupId", "DropQuery", "Location", "StatusCode", "ResponseMessage", "ContentType", "Condition"}
)

func AlbRuleHttpRuleCmd() *core.Command {
	albRuleHttpRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "httprule",
			Aliases:          []string{"http"},
			Short:            "Application Load Balancer Forwarding Rule Http Rule Operations",
			Long:             "The sub-commands of `ionosctl compute alb rule httprule` allow you to add, list, update, remove Application Load Balancer Forwarding Rule Http Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := albRuleHttpRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultAlbRuleHttpRuleCols, tabheaders.ColsMessage(allAlbRuleHttpRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(albRuleHttpRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = albRuleHttpRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAlbRuleHttpRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleListCmd())
	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleAddCmd())
	albRuleHttpRuleCmd.AddCommand(AlbRuleHttpRuleRemoveCmd())

	return core.WithConfigOverride(albRuleHttpRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
