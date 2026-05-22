package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	monitoring "github.com/ionos-cloud/sdk-go-bundle/products/monitoring/v2"
	"github.com/spf13/viper"
)

var pipelineCompleterCols = []table.Column{
	{Name: "Id", JSONPath: "id"},
	{Name: "Name", JSONPath: "properties.name"},
	{Name: "Status", JSONPath: "metadata.status"},
	{Name: "GrafanaEndpoint", JSONPath: "metadata.grafanaEndpoint"},
	{Name: "HttpEndpoint", JSONPath: "metadata.httpEndpoint"},
}

func PipelineIDs() []string {
	monClient := monitoring.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
	pipelines, _, err := monClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	t := table.New(pipelineCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(pipelines); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").AddInfo("GrafanaEndpoint").AddInfo("HttpEndpoint").AddInfo("Status").ToString()
}
