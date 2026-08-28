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
			Use:   "monitoring",
			Short: "Manage IONOS CLOUD Monitoring pipelines, ingest keys, and central monitoring",
			Long: `The Monitoring Service ingests, processes, and stores time-series metrics from your applications and infrastructure, and exposes them through a managed Grafana instance for visualization and alerting.

A pipeline exposes an HTTP ingest endpoint that agents push metrics to plus a managed Grafana endpoint for querying; an ingest key authenticates those pushes; and central monitoring lets other IONOS products forward their metrics to your pipelines automatically.

The service is regional: pipelines live in a specific location (e.g. de/txl) and every command targets one region via --location. An account may hold up to 10 pipelines by default (adjustable via Support).`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(pipeline.PipelineCommand())
	cmd.AddCommand(key.KeyCommand())
	cmd.AddCommand(central.CentralCommand())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Monitoring}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations)
}
