package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
			Example:   "ionosctl logging-service pipeline update --pipeline-id ID --json-properties PATH_TO_FILE",
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline you want to delete", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func runUpdateCmd(c *core.CommandConfig) error {
	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineId))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pipelineId,
	).Pipeline(
		pipelineToUpdate,
	).Execute()
	if err != nil {
		return err
	}

	return handlePipelinePrint(pipeline, c)
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.FlagPipelineId, constants.FlagJsonProperties},
		[]string{constants.FlagJsonPropertiesExample},
	)
}
