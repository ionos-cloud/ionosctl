package flowlog

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allFlowLogCols = []table.Column{
	{Name: "FlowLogId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Action", JSONPath: "properties.action", Default: true},
	{Name: "Direction", JSONPath: "properties.direction", Default: true},
	{Name: "Bucket", JSONPath: "properties.bucket", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func ApplicationLoadBalancerFlowLogCmd() *core.Command {
	applicationloadbalancerFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "Application Load Balancer FlowLog Operations",
			Long:             "The sub-commands of `ionosctl compute applicationloadbalancer flowlog` allow you to create, list, get, update, delete Application Load Balancer FlowLogs.",
			TraverseChildren: true,
		},
	}

	applicationloadbalancerFlowLogCmd.AddCommand(ApplicationLoadBalancerFlowLogListCmd())
	applicationloadbalancerFlowLogCmd.AddCommand(ApplicationLoadBalancerFlowLogGetCmd())
	applicationloadbalancerFlowLogCmd.AddCommand(ApplicationLoadBalancerFlowLogCreateCmd())
	applicationloadbalancerFlowLogCmd.AddCommand(ApplicationLoadBalancerFlowLogUpdateCmd())
	applicationloadbalancerFlowLogCmd.AddCommand(ApplicationLoadBalancerFlowLogDeleteCmd())

	return core.WithConfigOverride(applicationloadbalancerFlowLogCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
