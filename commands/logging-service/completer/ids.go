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
}

func LoggingServicePipelineIds() []string {
	pipelines, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
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
	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
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
