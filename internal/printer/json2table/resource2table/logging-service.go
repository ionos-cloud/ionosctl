package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
)

func ConvertLoggingServicePipelineLogToTable(log logging.PipelineResponse) ([]map[string]interface{}, error) {
	dests, ok := log.GetDestinationsOk()
	if !ok || dests == nil {
		return nil, fmt.Errorf("could not retrive Logging Service Pipeline Logs destination")
	}

	destinationsStrings := functional.Map(
		dests, func(dest logging.Destination) interface{} {
			return fmt.Sprintf("%v (%v days)", *dest.Type, *dest.RetentionInDays)
		},
	)

	logConverted, err := json2table.ConvertJSONToTable("", jsonpaths.LoggingServiceLogs, log)
	if err != nil {
		return nil, err
	}

	logConverted[0]["Destinations"] = destinationsStrings
	return logConverted, nil
}

func ConvertLoggingServicePipelineLogsToTable(pipeline logging.Pipeline) ([]map[string]interface{}, error) {
	logs, ok := pipeline.Properties.GetLogsOk()
	if !ok || logs == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline Logs")
	}

	var logsConverted []map[string]interface{}
	for _, log := range logs {
		dests, ok := log.GetDestinationsOk()
		if !ok || dests == nil {
			return nil, fmt.Errorf("could not retrive Logging Service Pipeline Logs destination")
		}

		destinationsStrings := functional.Map(
			dests, func(dest logging.Destination) interface{} {
				return fmt.Sprintf("%v (%v days)", *dest.Type, *dest.RetentionInDays)
			},
		)

		logConverted, err := json2table.ConvertJSONToTable("", jsonpaths.LoggingServiceLogs, log)
		if err != nil {
			return nil, err
		}

		logConverted[0]["Destinations"] = destinationsStrings

		logsConverted = append(logsConverted, logConverted...)
	}

	return logsConverted, nil
}

func ConvertLoggingServicePipelinesLogsToTable(pipelines logging.PipelineListResponse) (
	[]map[string]interface{}, error,
) {
	items, ok := pipelines.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline items")
	}

	var logsConverted []map[string]interface{}

	for _, pipeline := range items {
		logs, err := ConvertLoggingServicePipelineLogsToTable(pipeline)
		if err != nil {
			return nil, err
		}

		for _, l := range logs {
			l["PipelineId"] = pipeline.GetId()
		}

		logsConverted = append(logsConverted, logs...)

	}

	return logsConverted, nil
}
