package target

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultRuleTargetCols = []string{"TargetIp", "TargetPort", "Weight", "Check", "CheckInterval", "Maintenance"}
)

func NlbRuleTargetCmd() *core.Command {
	nlbRuleTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Network Load Balancer Forwarding Rule Target Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer rule target` allow you to add, list, update, remove Network Load Balancer Forwarding Rule Targets.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleTargetCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultRuleTargetCols, tabheaders.ColsMessage(defaultRuleTargetCols))
	_ = viper.BindPFlag(core.GetFlagName(nlbRuleTargetCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = nlbRuleTargetCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultRuleTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	nlbRuleTargetCmd.AddCommand(NlbRuleTargetListCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetAddCmd())
	nlbRuleTargetCmd.AddCommand(NlbRuleTargetRemoveCmd())

	return core.WithConfigOverride(nlbRuleTargetCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
