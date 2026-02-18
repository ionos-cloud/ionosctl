package flowlog

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultFlowLogCols = []string{"FlowLogId", "Name", "Action", "Direction", "Bucket", "State"}
)

func NetworkloadbalancerFlowLogCmd() *core.Command {
	networkloadbalancerFlowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"f", "fl"},
			Short:            "Network Load Balancer FlowLog Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer flowlog` allow you to create, list, get, update, delete Network Load Balancer FlowLogs.",
			TraverseChildren: true,
		},
	}

	networkloadbalancerFlowLogCmd.AddCommand(NetworkLoadBalancerFlowLogListCmd())
	networkloadbalancerFlowLogCmd.AddCommand(NetworkLoadBalancerFlowLogGetCmd())
	networkloadbalancerFlowLogCmd.AddCommand(NetworkLoadBalancerFlowLogCreateCmd())
	networkloadbalancerFlowLogCmd.AddCommand(NetworkLoadBalancerFlowLogUpdateCmd())
	networkloadbalancerFlowLogCmd.AddCommand(NetworkLoadBalancerFlowLogDeleteCmd())

	return core.WithConfigOverride(networkloadbalancerFlowLogCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
