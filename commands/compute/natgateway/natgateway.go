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
			Use:              "natgateway",
			Aliases:          []string{"nat", "ng"},
			Short:            "NAT Gateway Operations",
			Long:             "The sub-commands of `ionosctl compute natgateway` allow you to create, list, get, update, delete NAT Gateways.",
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
