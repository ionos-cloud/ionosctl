package rule

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "NatGatewayRuleId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Type", JSONPath: "properties.type"},
	{Name: "Protocol", JSONPath: "properties.protocol", Default: true},
	{Name: "SourceSubnet", JSONPath: "properties.sourceSubnet", Default: true},
	{Name: "PublicIp", JSONPath: "properties.publicIp", Default: true},
	{Name: "TargetSubnet", JSONPath: "properties.targetSubnet", Default: true},
	{Name: "TargetPortRangeStart", JSONPath: "properties.targetPortRange.start"},
	{Name: "TargetPortRangeEnd", JSONPath: "properties.targetPortRange.end"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func NatgatewayRuleCmd() *core.Command {
	natgatewayRuleCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "rule",
			Aliases: []string{"r"},
			Short:   "NAT Gateway Rule Operations",
			Long: `A NAT Gateway Rule is a source-NAT (SNAT) rule that decides which outbound traffic the gateway translates and which public IP it masquerades that traffic behind. Until at least one rule exists, a NAT Gateway forwards nothing.

Each rule matches a packet by:
  - source subnet  (--source-subnet, required)  the private source address range the rule applies to, e.g. the LAN's CIDR.
  - target subnet  (--target-subnet, optional)   restrict matching by destination address; omit to match any destination.
  - protocol       (--protocol)                  TCP, UDP, ICMP or ALL.
  - target port range (--port-range-start/--port-range-end) restrict matching by destination port; only valid for TCP/UDP.

Matched packets have their source address rewritten to the rule's public IP (--ip), which must be one of the public IPs already assigned to the parent NAT Gateway.`,
			TraverseChildren: true,
		},
	}
	natgatewayRuleCmd.AddColsFlag(allCols)

	natgatewayRuleCmd.AddCommand(NatgatewayRuleListCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleGetCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleCreateCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleUpdateCmd())
	natgatewayRuleCmd.AddCommand(NatgatewayRuleDeleteCmd())

	return core.WithConfigOverride(natgatewayRuleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
