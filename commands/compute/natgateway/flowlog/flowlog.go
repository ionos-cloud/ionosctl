package flowlog

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "FlowLogId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Action", JSONPath: "properties.action", Default: true},
	{Name: "Direction", JSONPath: "properties.direction", Default: true},
	{Name: "Bucket", JSONPath: "properties.bucket", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func NatgatewayFlowLogCmd() *core.Command {
	natgatewayFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "NAT Gateway FlowLog Operations",
			Long:             "The sub-commands of `ionosctl compute natgateway flowlog` allow you to create, list, get, update, delete NAT Gateway FlowLogs.",
			TraverseChildren: true,
		},
	}

	natgatewayFlowLogCmd.AddCommand(NatgatewayFlowLogListCmd())
	natgatewayFlowLogCmd.AddCommand(NatgatewayFlowLogGetCmd())
	natgatewayFlowLogCmd.AddCommand(NatgatewayFlowLogCreateCmd())
	natgatewayFlowLogCmd.AddCommand(NatgatewayFlowLogUpdateCmd())
	natgatewayFlowLogCmd.AddCommand(NatgatewayFlowLogDeleteCmd())

	return core.WithConfigOverride(natgatewayFlowLogCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
