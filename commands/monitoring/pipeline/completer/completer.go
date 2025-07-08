package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
)

func PipelineIDs() []string {
	pipeline, _, err := client.Must().Monitoring.PipelinesApi.PipelinesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	pipelineConverted, err := json2table.ConvertJSONToTable("", jsonpaths.MonitoringPipeline, pipeline)
	if err != nil {
		return nil
	}
	return completions.NewCompleter(pipelineConverted, "Id").AddInfo("Name").AddInfo("Status").AddInfo("GrafanaEndpoint").AddInfo("HttpEndpoint").ToString()
}
