package logs

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

func LogsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "logs",
			Short:            "Mongo Logs Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(LogsListCmd())

	return cmd
}

var allCols = []table.Column{
	{Name: "Instance", JSONPath: "Instance", Default: true},
	{Name: "Name", JSONPath: "Name", Default: true},
	{Name: "MessageNumber", JSONPath: "MessageNumber", Default: true},
	{Name: "Message", JSONPath: "Message"},
	{Name: "Time", JSONPath: "Time", Default: true},
}
