package logging_service

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func LoggingServiceCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "logging-service",
			Aliases:          []string{"log-svc"},
			Short:            "LaaS Operations. Manage and centralize your application/infrastructure's logs",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(pipeline.PipelineCmd())
	cmd.AddCommand(logs.LogsCmd())
	return cmd
}
