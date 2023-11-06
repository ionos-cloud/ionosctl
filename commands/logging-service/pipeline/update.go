package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
)

var (
	pipelineToUpdate            ionoscloud.CreateRequest
	jsonUpdatePropertiesExample = "{\n  \"properties\": {\n    \"name\": \"examplePipe\"," +
		"\n    \"logs\": [\n      {\n        \"source\": \"docker\",\n        \"tag\": \"tag1\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label1\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      },\n      {\n        \"source\": \"systemd\",\n        \"tag\": \"tag2\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label2\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 14\n          }\n        ]\n      },\n      {\n        \"source\": \"kubernetes\",\n        \"tag\": \"tag3\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label3\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 30\n          }\n        ]\n      },\n      {\n        \"source\": \"generic\",\n        \"tag\": \"tag4\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label4\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      }\n    ]\n  }\n}\n"
)

func PipelineUpdateCmd() *core.Command {
	cmd := core.NewCommandWithJsonProperties(
		context.Background(), nil, jsonUpdatePropertiesExample, &pipelineToUpdate, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "update",
			ShortDesc: "Create a logging pipeline",
			Example:   "ionosctl logging-service pipeline create --json-properties PATH_TO_FILE",
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))

	return cmd
}

func runUpdateCmd(c *core.CommandConfig) error {
	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPost(context.Background()).Pipeline(
		pipelineToCreate,
	).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.LoggingServicePipeline, pipeline,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.FlagJsonProperties}, []string{constants.FlagJsonPropertiesExample},
	)
}
