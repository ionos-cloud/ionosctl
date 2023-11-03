package logs

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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

var (
	allCols     = []string{"Instance", "Name", "MessageNumber", "Message", "Time"}
	defaultCols = []string{"Instance", "Name", "MessageNumber", "Time"}
)
