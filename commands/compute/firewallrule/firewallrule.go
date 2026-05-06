package firewallrule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allFirewallRuleCols = []table.Column{
	{Name: "FirewallRuleId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Protocol", JSONPath: "properties.protocol", Default: true},
	{Name: "PortRangeStart", JSONPath: "properties.portRangeStart", Default: true},
	{Name: "PortRangeEnd", JSONPath: "properties.portRangeEnd", Default: true},
	{Name: "Direction", JSONPath: "properties.type", Default: true},
	{Name: "IPVersion", JSONPath: "properties.ipVersion", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "SourceMac", JSONPath: "properties.sourceMac"},
	{Name: "SourceIP", JSONPath: "properties.sourceIp"},
	{Name: "DestinationIP", JSONPath: "properties.targetIp"},
	{Name: "IcmpCode", JSONPath: "properties.icmpCode"},
	{Name: "IcmpType", JSONPath: "properties.icmpType"},
}

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
	firewallRuleCmd.AddColsFlag(allFirewallRuleCols)

	firewallRuleCmd.AddCommand(FirewallRuleListCmd())
	firewallRuleCmd.AddCommand(FirewallRuleGetCmd())
	firewallRuleCmd.AddCommand(FirewallRuleCreateCmd())
	firewallRuleCmd.AddCommand(FirewallRuleUpdateCmd())
	firewallRuleCmd.AddCommand(FirewallRuleDeleteCmd())

	return core.WithConfigOverride(firewallRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
