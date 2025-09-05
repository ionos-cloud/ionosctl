package logging_service

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/central"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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
	cmd.AddCommand(central.CentralCommand())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Logging}, constants.LoggingApiRegionalURL, constants.LoggingLocations)
}
