package logging_service

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "logging-service",
			Aliases: []string{"log-svc"},
			Short: "Logging Service Operations. " +
				"Manage and centralize your application/infrastructure's logs",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(pipeline.PipelineCmd())
	cmd.AddCommand(logs.LogsCmd())

	return core.WithRegionalFlags(cmd, constants.LoggingApiRegionalURL, constants.LoggingLocations)
}
