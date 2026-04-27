package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
)

var pipelineCompleterCols = []table.Column{
	{Name: "Id", JSONPath: "id"},
	{Name: "Name", JSONPath: "properties.name"},
	{Name: "Status", JSONPath: "metadata.status"},
	{Name: "GrafanaEndpoint", JSONPath: "metadata.grafanaEndpoint"},
	{Name: "HttpEndpoint", JSONPath: "metadata.httpEndpoint"},
}

func PipelineIDs() []string {
	pipelines, _, err := client.Must().Monitoring.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	t := table.New(pipelineCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(pipelines); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").AddInfo("GrafanaEndpoint").AddInfo("HttpEndpoint").AddInfo("Status").ToString()
}
