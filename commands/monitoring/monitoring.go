package monitoring

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "monitoring",
			Short:            "Monitoring Service is a cloud-based service that allows you to ingest, aggregate, and analyze data to enhance your understanding of your system's performance and behavior",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(pipeline.PipelineCommand())

	return core.WithRegionalFlags(cmd, constants.MonitoringApiRegionalURL, constants.MonitoringLocations)
}
