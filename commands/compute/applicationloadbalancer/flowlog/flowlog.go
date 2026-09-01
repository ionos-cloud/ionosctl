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
			Use:     "flowlog",
			Aliases: []string{"f", "fl"},
			Short:   "Application Load Balancer FlowLog Operations",
			Long: `A flow log captures metadata about the traffic handled by an Application Load Balancer and streams it to an IONOS Object Storage (S3) bucket for auditing, troubleshooting and compliance. Flow logs record connection metadata (not payloads).

Each flow log is scoped by:
  * --action    which connections to log: ACCEPTED (allowed), REJECTED (denied), or ALL.
  * --direction which traffic to log relative to the ALB: INGRESS (inbound), EGRESS (outbound), or BIDIRECTIONAL.
  * --s3bucket  the name of an existing IONOS Object Storage bucket that receives the log records.

The sub-commands below let you create, list, get, update and delete flow logs.`,
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
