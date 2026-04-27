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
			Use:              "pipeline",
			Aliases:          []string{"p", "pipe"},
			Short:            "A metric pipeline refers to an instance or configuration of the Monitoring Service ",
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
