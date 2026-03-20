package targetgroup

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTargetGroupCols = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "State"}
	allTargetGroupCols     = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "Retries",
		"Path", "Method", "MatchType", "Response", "Regex", "Negate", "State"}
)

func TargetGroupCmd() *core.Command {
	targetGroupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "targetgroup",
			Aliases:          []string{"tg"},
			Short:            "Target Group Operations",
			Long:             "The sub-commands of `ionosctl compute targetgroup` allow you to see information, to create, update, delete Target Groups.",
			TraverseChildren: true,
		},
	}
	globalFlags := targetGroupCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTargetGroupCols, tabheaders.ColsMessage(allTargetGroupCols))
	_ = viper.BindPFlag(core.GetFlagName(targetGroupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = targetGroupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTargetGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	targetGroupCmd.AddCommand(TargetGroupListCmd())
	targetGroupCmd.AddCommand(TargetGroupGetCmd())
	targetGroupCmd.AddCommand(TargetGroupCreateCmd())
	targetGroupCmd.AddCommand(TargetGroupUpdateCmd())
	targetGroupCmd.AddCommand(TargetGroupDeleteCmd())
	targetGroupCmd.AddCommand(TargetGroupTargetCmd())

	return core.WithConfigOverride(targetGroupCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
