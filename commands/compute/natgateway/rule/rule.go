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
			Use:              "rule",
			Aliases:          []string{"r"},
			Short:            "NAT Gateway Rule Operations",
			Long:             "The sub-commands of `ionosctl compute natgateway rule` allow you to create, list, get, update, delete NAT Gateway Rules.",
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
