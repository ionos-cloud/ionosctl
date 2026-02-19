package firewallrule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultFirewallRuleCols = []string{"FirewallRuleId", "Name", "Protocol", "PortRangeStart", "PortRangeEnd", "Direction", "IPVersion", "State"}
	allFirewallRuleCols     = []string{"FirewallRuleId", "Name", "Protocol", "SourceMac", "SourceIP", "DestinationIP", "PortRangeStart", "PortRangeEnd",
		"IcmpCode", "IcmpType", "Direction", "IPVersion", "State"}
)

func FirewallruleCmd() *core.Command {
	firewallRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "firewallrule",
			Aliases:          []string{"f", "fr", "firewall"},
			Short:            "Firewall Rule Operations",
			Long:             "The sub-commands of `ionosctl compute firewallrule` allow you to create, list, get, update, delete Firewall Rules.",
			TraverseChildren: true,
		},
	}
	globalFlags := firewallRuleCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultFirewallRuleCols, tabheaders.ColsMessage(allFirewallRuleCols))
	_ = viper.BindPFlag(core.GetFlagName(firewallRuleCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = firewallRuleCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allFirewallRuleCols, cobra.ShellCompDirectiveNoFileComp
	})

	firewallRuleCmd.AddCommand(FirewallRuleListCmd())
	firewallRuleCmd.AddCommand(FirewallRuleGetCmd())
	firewallRuleCmd.AddCommand(FirewallRuleCreateCmd())
	firewallRuleCmd.AddCommand(FirewallRuleUpdateCmd())
	firewallRuleCmd.AddCommand(FirewallRuleDeleteCmd())

	return core.WithConfigOverride(firewallRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
