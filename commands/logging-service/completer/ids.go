package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	logging "github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

var pipelineCompleterCols = []table.Column{
	{Name: "Id", JSONPath: "id"},
	{Name: "Name", JSONPath: "properties.name"},
}

func LoggingServicePipelineIds() []string {
	logClient := logging.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
	pipelines, _, err := logClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	t := table.New(pipelineCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(pipelines); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").ToString()
}

func LoggingServiceLogTags(pipelineId string) []string {
	logClient := logging.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))
	pipeline, _, err := logClient.PipelinesApi.PipelinesFindById(
		context.Background(),
		pipelineId,
	).Execute()
	if err != nil {
		return nil
	}

	var rows []map[string]any
	for _, log := range pipeline.Properties.Logs {
		rows = append(rows, map[string]any{
			"Tag":      log.Tag,
			"Source":   log.Source,
			"Protocol": log.Protocol,
		})
	}

	return completions.NewCompleter(rows, "Tag").AddInfo("Source").AddInfo("Protocol", "(%v)").ToString()
}
