package pipeline

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "GrafanaEndpoint", JSONPath: "metadata.grafanaEndpoint", Default: true},
	{Name: "HttpEndpoint", JSONPath: "metadata.httpEndpoint", Default: true},
	{Name: "Status", JSONPath: "metadata.status", Default: true},
}

func PipelineCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "pipeline",
			Aliases: []string{"p", "pipe"},
			Short:   "Create and manage monitoring pipelines (metric ingest + Grafana endpoints)",
			Long: `A pipeline is the core Monitoring Service resource: it receives, processes, and stores metrics from one or more sources. Creating a pipeline returns two endpoints and an ingest key:
  * HttpEndpoint    - the HTTP endpoint agents push metrics to, authenticating with the pipeline's ingest key.
  * GrafanaEndpoint - the managed Grafana base URL where the ingested metrics can be queried and turned into dashboards.
  * ingest key      - shown only once, at creation. It authenticates every metric push; store it securely. Use 'monitoring key create' to rotate it (which invalidates the old key).

Pipelines are regional: each lives in the location it was created in. Commands that act on a single pipeline require --location; 'pipeline list' searches all regions by default. An account may hold up to 10 pipelines by default.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)
	cmd.AddCommand(MonitoringListCmd())
	cmd.AddCommand(MonitoringFindByIdCmd())
	cmd.AddCommand(MonitoringDeleteCmd())
	cmd.AddCommand(MonitoringPostCmd())
	cmd.AddCommand(MonitoringPutCmd())

	return cmd
}
