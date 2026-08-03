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
			Use:     "flowlog",
			Aliases: []string{"f", "fl"},
			Short:   "NAT Gateway FlowLog Operations",
			Long: `A NAT Gateway FlowLog records the traffic that passes through the gateway and ships those records to an IONOS Object Storage (S3) bucket you own, for auditing, troubleshooting and compliance.

Each flowlog is scoped by two patterns:
  - action    (--action)    which flows to log by outcome: ACCEPTED, REJECTED, or ALL.
  - direction (--direction) which flows to log by direction relative to the gateway: INGRESS, EGRESS, or BIDIRECTIONAL.

Logs are delivered to the bucket named by --s3bucket, which must already exist in your IONOS Object Storage. The bucket is not created for you.`,
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
