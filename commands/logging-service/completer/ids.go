package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
)

func LoggingServicePipelineIds() []string {
	pipelines, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	pipelinesConverted, err := json2table.ConvertJSONToTable(
		"items", jsonpaths.LoggingServicePipeline, pipelines,
	)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(pipelinesConverted, "Id").AddInfo("Name").ToString()
}

func LoggingServiceLogTags(pipelineId string) []string {
	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(),
		pipelineId,
	).Execute()
	if err != nil {
		return nil
	}

	logsConverted, err := resource2table.ConvertLoggingServicePipelineLogsToTable(pipeline)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(logsConverted, "Tag").AddInfo("Source").AddInfo("Protocol", "(%v)").ToString()
}
