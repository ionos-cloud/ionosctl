package applicationloadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/applicationloadbalancer/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultApplicationLoadBalancerCols = []string{"ApplicationLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
	allApplicationLoadBalancerCols     = []string{"ApplicationLoadBalancerId", "DatacenterId", "Name", "ListenerLan", "Ips", "TargetLan", "PrivateIps", "State"}
)

func ApplicationLoadBalancerCmd() *core.Command {
	applicationloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "applicationloadbalancer",
			Aliases:          []string{"alb"},
			Short:            "Application Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl compute applicationloadbalancer` allow you to create, list, get, update, delete Application Load Balancers.",
			TraverseChildren: true,
		},
	}
	globalFlags := applicationloadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultApplicationLoadBalancerCols, tabheaders.ColsMessage(allApplicationLoadBalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(applicationloadbalancerCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = applicationloadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allApplicationLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerListCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerGetCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerCreateCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerUpdateCmd())
	applicationloadbalancerCmd.AddCommand(ApplicationLoadBalancerDeleteCmd())

	applicationloadbalancerCmd.AddCommand(rule.ApplicationLoadBalancerRuleCmd())
	applicationloadbalancerCmd.AddCommand(flowlog.ApplicationLoadBalancerFlowLogCmd())

	return core.WithConfigOverride(applicationloadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
