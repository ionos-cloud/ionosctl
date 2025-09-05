package monitoring

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/central"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/key"
	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
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
	cmd.AddCommand(central.CentralCommand())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Monitoring}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations)
}
