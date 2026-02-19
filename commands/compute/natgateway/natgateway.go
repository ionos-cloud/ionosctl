package natgateway

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/lan"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultNatGatewayCols = []string{"NatGatewayId", "Name", "PublicIps", "State"}
	allNatGatewayCols     = []string{"NatGatewayId", "Name", "PublicIps", "State", "DatacenterId"}
)

func NatgatewayCmd() *core.Command {
	natgatewayCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "natgateway",
			Aliases:          []string{"nat", "ng"},
			Short:            "NAT Gateway Operations",
			Long:             "The sub-commands of `ionosctl compute natgateway` allow you to create, list, get, update, delete NAT Gateways.",
			TraverseChildren: true,
		},
	}
	globalFlags := natgatewayCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNatGatewayCols, tabheaders.ColsMessage(defaultNatGatewayCols))
	_ = viper.BindPFlag(core.GetFlagName(natgatewayCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = natgatewayCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultNatGatewayCols, cobra.ShellCompDirectiveNoFileComp
	})

	natgatewayCmd.AddCommand(NatgatewayListCmd())
	natgatewayCmd.AddCommand(NatgatewayGetCmd())
	natgatewayCmd.AddCommand(NatgatewayCreateCmd())
	natgatewayCmd.AddCommand(NatgatewayUpdateCmd())
	natgatewayCmd.AddCommand(NatgatewayDeleteCmd())

	natgatewayCmd.AddCommand(lan.NatgatewayLanCmd())
	natgatewayCmd.AddCommand(rule.NatgatewayRuleCmd())
	natgatewayCmd.AddCommand(flowlog.NatgatewayFlowLogCmd())

	return core.WithConfigOverride(natgatewayCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
