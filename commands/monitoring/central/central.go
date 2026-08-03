package central

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Enabled", JSONPath: "properties.enabled", Default: true},
	{Name: "GrafanaEndpoint", JSONPath: "metadata.grafanaEndpoint", Default: true},
	{Name: "Products", JSONPath: "metadata.products.*", Default: true},
}

func CentralCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "central",
			Aliases: []string{"c"},
			Short:   "View and toggle central monitoring for a region",
			Long: `Central monitoring is a per-region toggle. When enabled, other IONOS CLOUD products in that region forward their metrics to the Monitoring Service on your behalf, without you having to configure a push agent for each one. When disabled, only metrics you push explicitly (via a pipeline's ingest key) are collected.

The state applies to the whole region selected with --location, not to an individual pipeline. 'central get' reports whether it is enabled, the Grafana endpoint, and which products are currently forwarding metrics.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(CentralFindByIdCmd())
	cmd.AddCommand(CentralDisable())
	cmd.AddCommand(CentralEnable())

	return cmd
}

func enable(c *core.CommandConfig, enabled bool) error {

	input := monitoring.CentralMonitoring{Enabled: enabled}

	r, _, err := client.Must().Monitoring.CentralApi.CentralPut(context.Background(), "").
		CentralMonitoringEnsure(monitoring.CentralMonitoringEnsure{
			Properties: input,
		}).Execute()
	if err != nil {
		return fmt.Errorf("failed changing the enabled state: %w", err)
	}

	return c.Printer(allCols).Print(r)
}
