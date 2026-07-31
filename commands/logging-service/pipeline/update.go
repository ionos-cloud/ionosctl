package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

var (
	pipelineToUpdate            logging.PipelinePatch
	jsonUpdatePropertiesExample = "{\n  \"properties\": {\n    \"name\": \"examplePipe\"," +
		"\n    \"logs\": [\n      {\n        \"source\": \"docker\",\n        \"tag\": \"tag1\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label1\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      },\n      {\n        \"source\": \"systemd\",\n        \"tag\": \"tag2\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label2\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 14\n          }\n        ]\n      },\n      {\n        \"source\": \"kubernetes\",\n        \"tag\": \"tag3\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label3\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 30\n          }\n        ]\n      },\n      {\n        \"source\": \"generic\",\n        \"tag\": \"tag4\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label4\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      }\n    ]\n  }\n}\n"
)

func PipelineUpdateCmd() *core.Command {
	cmd := core.NewCommandWithJsonProperties(
		context.Background(), nil, jsonUpdatePropertiesExample, &pipelineToUpdate, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "update",
			ShortDesc: "Update a logging pipeline",
			LongDesc: `Replace a pipeline's properties (its name and the full set of logs) from a JSON file. The JSON is sent as a patch, so the 'logs' array you provide REPLACES the existing logs - include every log you want to keep, not just the changed ones. Use --json-properties-example to print a template with the expected shape.

To tweak a single log stream in place instead, prefer 'ionosctl logging-service logs update', which edits one log by its tag without rewriting the whole pipeline.

The pipeline's ingestion/Grafana addresses and its key are unaffected by an update.`,
			Example: `# Replace the pipeline's name and logs from a file
ionosctl logging-service pipeline update --location de/txl --pipeline-id ID --json-properties ./pipeline.json

# Print a JSON template to edit
ionosctl logging-service pipeline update --json-properties-example`,
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline to update", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func runUpdateCmd(c *core.CommandConfig) error {
	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pipelineId,
	).PipelinePatch(
		pipelineToUpdate,
	).Execute()
	if err != nil {
		return err
	}

	return handlePipelinePrint(pipeline, c)
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsSetsAndLocation(
		[]string{constants.FlagLoggingPipelineId, constants.FlagJsonProperties},
		[]string{constants.FlagJsonPropertiesExample},
	)
}
