package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pipelineToUpdate            ionoscloud.PipelinePatch
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
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServicePipelineIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func runUpdateCmd(c *core.CommandConfig) error {
	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))

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
		c.Command, c.NS, []string{constants.FlagLoggingPipelineId, constants.FlagJsonProperties},
		[]string{constants.FlagJsonPropertiesExample},
	)
}
