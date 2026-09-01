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
			Use:     "firewallrule",
			Aliases: []string{"f", "fr", "firewall"},
			Short:   "Firewall Rule Operations",
			Long: `Manage Firewall Rules attached to a NIC (Network Interface Card).

A Firewall Rule filters traffic on one specific NIC, identified by the tuple --datacenter-id / --server-id / --nic-id. Rules therefore live per-NIC: two NICs on the same server each have their own independent rule set, and a rule created on one NIC does not affect the other.

The NIC firewall is default-deny. When a NIC's firewall is active but has no rules, ALL traffic is blocked. Each rule you add whitelists (allows) a specific slice of traffic; you build up the allowed set rule by rule. Activate the firewall on the NIC itself (see 'ionosctl compute nic update --firewall-active') - firewall rules are only enforced while the NIC's firewall is active.

Each rule has a direction (--direction): INGRESS filters traffic coming from outside toward the NIC, EGRESS filters traffic leaving the NIC toward outside. A rule matches on any combination of protocol, source MAC, source IP, target IP, and - depending on protocol - a port range (TCP/UDP) or an ICMP type/code. Unset match fields act as wildcards (allow all).`,
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
