package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func PipelineGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "get",
			ShortDesc: "Retrieve a logging pipeline by ID",
			LongDesc:  `Show a single pipeline, including its TCP/HTTP ingestion addresses and Grafana address. The pipeline key is not returned here; generate one with 'pipeline key'.`,
			Example:   "ionosctl logging-service pipeline get --location de/txl --pipeline-id ID",
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline to retrieve", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func preRunGetCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagLoggingPipelineId)
}

func runGetCmd(c *core.CommandConfig) error {
	pipelineId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pipelineId,
	).Execute()
	if err != nil {
		return err
	}

	return handlePipelinePrint(pipeline, c)
}
