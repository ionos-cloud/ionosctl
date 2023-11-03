package logs

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Source", "Protocol", "Public", "Destinations"}
	allCols     = []string{"Source", "Protocol", "Public", "Destinations", "Tag", "Labels", "PipelineId"}
)

func LogsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use: "logs",
			Short: "The subcommands of `ionosctl logging-service logs` allow you to manage logging pipelines. " +
				"They are the backbone of a centralized logging system, " +
				"referring to an instance or configuration of the logging service you can create",
		},
	}

	cmd.AddCommand(LogsListCmd())
	return cmd
}
