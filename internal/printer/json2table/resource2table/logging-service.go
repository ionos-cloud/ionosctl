package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
)

func ConvertLoggingServicePipelineLogToTable(log ionoscloud.ResponsePipeline) ([]map[string]interface{}, error) {
	destinations, ok := log.GetDestinationsOk()
	if !ok || destinations == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline Log destinations")
	}

	var destinationsStrings []interface{}
	for _, dest := range *destinations {
		destinationsStrings = append(
			destinationsStrings, fmt.Sprintf(
				"%v (%v days)", *dest.Type,
				*dest.RetentionInDays,
			),
		)
	}

	logConverted, err := json2table.ConvertJSONToTable("", jsonpaths.LoggingServiceLogs, log)
	if err != nil {
		return nil, err
	}

	logConverted[0]["Destinations"] = destinationsStrings

	return logConverted, nil
}

func ConvertLoggingServicePipelineLogsToTable(pipelines ionoscloud.PipelineListResponse) (
	[]map[string]interface{}, error,
) {
	items, ok := pipelines.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline items")
	}

	var logsConverted []map[string]interface{}

	for _, pipeline := range *items {
		for _, log := range *pipeline.Properties.Logs {
			l, err := ConvertLoggingServicePipelineLogToTable(log)
			if err != nil {
				return nil, err
			}

			l[0]["PipelineId"] = *pipeline.Id
			logsConverted = append(logsConverted, l...)
		}
	}

	return logsConverted, nil
}
