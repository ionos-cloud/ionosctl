package networkloadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/networkloadbalancer/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultNetworkLoadBalancerCols = []string{"NetworkLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "LbPrivateIps", "State"}
	allNetworkLoadBalancerCols     = []string{"NetworkLoadBalancerId", "Name", "ListenerLan", "Ips", "TargetLan", "LbPrivateIps", "State", "DatacenterId"}
)

func NetworkloadbalancerCmd() *core.Command {
	networkloadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "networkloadbalancer",
			Aliases:          []string{"nlb"},
			Short:            "Network Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer` allow you to create, list, get, update, delete Network Load Balancers.",
			TraverseChildren: true,
		},
	}
	globalFlags := networkloadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNetworkLoadBalancerCols, tabheaders.ColsMessage(defaultNetworkLoadBalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(networkloadbalancerCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = networkloadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNetworkLoadBalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerListCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerGetCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerCreateCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerUpdateCmd())
	networkloadbalancerCmd.AddCommand(NetworkLoadBalancerDeleteCmd())

	networkloadbalancerCmd.AddCommand(flowlog.NetworkloadbalancerFlowLogCmd())
	networkloadbalancerCmd.AddCommand(rule.NetworkloadbalancerRuleCmd())

	return core.WithConfigOverride(networkloadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
