package pipeline

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/cobra"
)

var allCols = []table.Column{
	{Name: "Id", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "GrafanaAddress", JSONPath: "properties.grafanaAddress", Default: true},
	{Name: "TCPAddress", JSONPath: "properties.tcpAddress"},
	{Name: "HTTPAddress", JSONPath: "properties.httpAddress"},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func PipelineCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:     "pipeline",
			Aliases: []string{"p", "pipelines"},
			Short:   "Manage logging pipelines",
			Long: `A pipeline is the top-level Logging Service resource. It provisions a managed ingestion tier plus a Loki/Grafana backend, and exposes:
  - a TCP address and an HTTP address that your log shippers push logs to,
  - a Grafana address where you query and visualise those logs.

Each pipeline holds one or more logs (see 'ionosctl logging-service logs'), where every log ties a source and tag to a protocol and a destination with its own retention. Pipelines are regional: create, get, update and key all act on a single --location.

Authentication to the ingestion endpoints is done with a pipeline key. The key is generated separately with 'pipeline key' and can be rotated at any time.`,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(PipelineListCmd())
	cmd.AddCommand(PipelineGetCmd())
	cmd.AddCommand(PipelineDeleteCmd())
	cmd.AddCommand(PipelineCreateCmd())
	cmd.AddCommand(PipelineUpdateCmd())
	cmd.AddCommand(PipelineKeyCmd())
	return cmd
}

func handlePipelinePrint(p logging.PipelineRead, c *core.CommandConfig) error {
	return c.Printer(allCols).Print(p)
}
