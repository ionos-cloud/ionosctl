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
			Use:              "flowlog",
			Aliases:          []string{"fl"},
			Short:            "FlowLog Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl compute flowlog` + "`" + ` allow you to create, list, get, delete FlowLogs on specific NICs.`,
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
