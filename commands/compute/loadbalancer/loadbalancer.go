package loadbalancer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/loadbalancer/nic"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultLoadbalancerCols = []string{"LoadBalancerId", "Name", "Dhcp", "State"}
	allLoadbalancerCols     = []string{"LoadBalancerId", "Name", "Dhcp", "State", "Ip", "DatacenterId"}
)

func LoadBalancerCmd() *core.Command {
	loadbalancerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "loadbalancer",
			Aliases:          []string{"lb"},
			Short:            "Load Balancer Operations",
			Long:             "The sub-commands of `ionosctl compute loadbalancer` manage your Load Balancers on your account. With Load Balancers you can distribute traffic between your servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := loadbalancerCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultLoadbalancerCols, tabheaders.ColsMessage(allLoadbalancerCols))
	_ = viper.BindPFlag(core.GetFlagName(loadbalancerCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = loadbalancerCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allLoadbalancerCols, cobra.ShellCompDirectiveNoFileComp
	})

	loadbalancerCmd.AddCommand(LoadBalancerListCmd())
	loadbalancerCmd.AddCommand(LoadBalancerGetCmd())
	loadbalancerCmd.AddCommand(LoadBalancerCreateCmd())
	loadbalancerCmd.AddCommand(LoadBalancerUpdateCmd())
	loadbalancerCmd.AddCommand(LoadBalancerDeleteCmd())
	loadbalancerCmd.AddCommand(nic.LoadBalancerNicCmd())

	return core.WithConfigOverride(loadbalancerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
