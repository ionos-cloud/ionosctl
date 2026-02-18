package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/rule/target"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultForwardingRuleCols = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State"}
	allForwardingRuleCols     = []string{"ForwardingRuleId", "Name", "Algorithm", "Protocol", "ListenerIp", "ListenerPort", "State",
		"ClientTimeout", "ConnectTimeout", "TargetTimeout", "Retries"}
)

func NetworkloadbalancerRuleCmd() *core.Command {
	nlbRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r", "forwardingrule"},
			Short:            "Network Load Balancer Forwarding Rule Operations",
			Long:             "The sub-commands of `ionosctl nlb rule` allow you to create, list, get, update, delete Network Load Balancer Forwarding Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultForwardingRuleCols, tabheaders.ColsMessage(allForwardingRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(nlbRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = nlbRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allForwardingRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleListCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleGetCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleCreateCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleUpdateCmd())
	nlbRuleCmd.AddCommand(NetworkLoadBalancerForwardingRuleDeleteCmd())

	nlbRuleCmd.AddCommand(target.NlbRuleTargetCmd())

	return core.WithConfigOverride(nlbRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
