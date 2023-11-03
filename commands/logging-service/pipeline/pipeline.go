package pipeline

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultCols = []string{"Id", "Name", "GrafanaAddress", "CreatedDate"}
	allCols     = []string{"Id", "Name", "GrafanaAddress", "TCPAddress", "HTTPAddress", "CreatedDate"}
)

func PipelineCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "pipeline",
			Aliases: []string{"p", "pipelines"},
			Short: "The subcommands of `ionosctl logging-service pipeline` allow you to manage logging pipelines. " +
				"They are the backbone of a centralized logging system, " +
				"referring to an instance or configuration of the logging service you can create",
		},
	}

	cmd.AddCommand(PipelineListCmd())
	return cmd
}
