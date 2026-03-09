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

	cmd.AddCommand(LogsListCmd())

	return cmd
}

var allCols = []table.Column{
	{Name: "Instance", Default: true},
	{Name: "Name", Default: true},
	{Name: "MessageNumber", Default: true},
	{Name: "Message"},
	{Name: "Time", JSONPath: "time", Default: true},
}
