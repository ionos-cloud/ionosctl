package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultNatGatewayRuleCols = []string{"NatGatewayRuleId", "Name", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "State"}
	allNatGatewayRuleCols     = []string{"NatGatewayRuleId", "Name", "Type", "Protocol", "SourceSubnet", "PublicIp", "TargetSubnet", "TargetPortRangeStart", "TargetPortRangeEnd", "State"}
)

func NatgatewayRuleCmd() *core.Command {
	natgatewayRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "rule",
			Aliases:          []string{"r"},
			Short:            "NAT Gateway Rule Operations",
			Long:             "The sub-commands of `ionosctl compute natgateway rule` allow you to create, list, get, update, delete NAT Gateway Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := natgatewayRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultNatGatewayRuleCols, tabheaders.ColsMessage(allNatGatewayRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(natgatewayRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = natgatewayRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allNatGatewayRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	natgatewayRuleCmd.AddCommand(NatgatewayRuleListCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleGetCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleCreateCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleUpdateCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleDeleteCmd())

	return core.WithConfigOverride(natgatewayRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
