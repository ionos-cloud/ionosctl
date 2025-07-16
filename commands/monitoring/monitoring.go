package monitoring

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/key"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "monitoring",
			Short:            "Monitoring is a cloud service that collects and analyzes data to improve system performance",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(pipeline.PipelineCommand())
	cmd.AddCommand(key.KeyCommand())

	return core.WithRegionalConfigOverride(cmd, "monitoring", constants.MonitoringApiRegionalURL, constants.MonitoringLocations)
}
