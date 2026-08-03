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

func FlowlogCmd() *core.Command {
	flowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "flowlog",
			Aliases: []string{"fl"},
			Short:   "FlowLog Operations",
			Long: `The sub-commands of ` + "`" + `ionosctl compute flowlog` + "`" + ` let you manage Flow Logs on a NIC.

A Flow Log captures metadata about the network traffic that passes through a NIC (the connection 5-tuple, whether it was accepted or rejected, byte/packet counts, etc.) and delivers those records to an IONOS Object Storage (S3) bucket for auditing, troubleshooting, and security analysis. Flow Logs record traffic metadata only, not packet payloads.

Every Flow Log is scoped to one NIC, so all of these commands require the full path to that NIC: --datacenter-id, --server-id and --nic-id. A Flow Log is defined by:
  * action     - which traffic to capture: ACCEPTED, REJECTED, or ALL
  * direction  - which flow direction to capture: INGRESS, EGRESS, or BIDIRECTIONAL
  * bucket     - the name of an EXISTING IONOS Object Storage bucket that receives the log files

NOTE: Disable/delete the Flow Log before deleting the Object Storage bucket it writes to, otherwise log delivery breaks.`,
			TraverseChildren: true,
		},
	}
	flowLogCmd.AddColsFlag(allCols)

	flowLogCmd.AddCommand(FlowLogListCmd())
	flowLogCmd.AddCommand(FlowLogGetCmd())
	flowLogCmd.AddCommand(FlowLogCreateCmd())
	flowLogCmd.AddCommand(FlowLogDeleteCmd())

	return core.WithConfigOverride(flowLogCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
