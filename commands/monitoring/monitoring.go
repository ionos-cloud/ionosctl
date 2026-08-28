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

The domain model has three parts:
  * pipeline - the core resource. Each pipeline exposes an HTTP ingest endpoint that agents push metrics to, and a Grafana endpoint for querying and dashboards.
  * key      - the per-pipeline ingest key (API key) that authenticates metric pushes. Rotating a key immediately invalidates the previous one.
  * central  - central monitoring, an account/region-level toggle that lets other IONOS products forward their metrics to your pipelines automatically.

The service is regional: pipelines live in a specific location (e.g. de/txl) and every command targets one region via --location. An account may hold up to 10 pipelines by default (adjustable via Support).`,
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(pipeline.PipelineCommand())
	cmd.AddCommand(key.KeyCommand())
	cmd.AddCommand(central.CentralCommand())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Monitoring}, constants.MonitoringApiRegionalURL, constants.MonitoringLocations)
}
