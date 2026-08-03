package natgateway

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/flowlog"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/lan"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/natgateway/rule"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "NatGatewayId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "PublicIps", JSONPath: "properties.publicIps", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "DatacenterId", JSONPath: "href"},
}

func NatgatewayCmd() *core.Command {
	natgatewayCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "natgateway",
			Aliases: []string{"nat", "ng"},
			Short:   "NAT Gateway Operations",
			Long: `A NAT Gateway gives servers on a private LAN outbound access to the public internet without assigning each server its own public IP. It does this via source-NAT (SNAT): as a packet leaves toward the internet, the gateway rewrites the packet's private source address to one of the gateway's own public IPs, then translates the return traffic back. Because connections can only be initiated from the inside out, private servers stay unreachable from the internet.

A NAT Gateway lives inside one Virtual Data Center (` + "`" + `--datacenter-id` + "`" + `) and owns one or more public IPs (` + "`" + `--ips` + "`" + `), which must be reserved IP addresses in the same location as the datacenter. The resource is made up of three sub-resources:

  - lan     Attach the gateway to a private LAN so servers on that LAN can route their outbound traffic through it. Each attached LAN also gets one or more gateway IPs (the next-hop address servers use).
  - rule    SNAT rules that decide which traffic is translated and which public IP is used. A packet is matched by source subnet (and optionally destination subnet, protocol and target port range), then masqueraded behind the rule's public IP. Without any rule, no traffic is translated.
  - flowlog Capture accepted/rejected traffic that passes the gateway into an IONOS Object Storage (S3) bucket for auditing and troubleshooting.

Typical flow: create the gateway with its public IPs, attach the private LAN(s), then add SNAT rule(s) so the servers on those LANs reach the internet.`,
			TraverseChildren: true,
		},
	}
	natgatewayCmd.AddColsFlag(allCols)

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
