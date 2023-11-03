package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
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
