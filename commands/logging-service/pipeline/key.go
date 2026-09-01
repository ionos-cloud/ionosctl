package pipeline

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func PipelineKeyCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "key",
			ShortDesc: "Generate (rotate) a logging pipeline's authentication key",
			LongDesc: `Generate a new key for a pipeline and print it. Log shippers present this key when pushing to the pipeline's TCP/HTTP ingestion address, so it acts as the pipeline's write credential.

Generating a key ROTATES it: any previous key is immediately invalidated, so update your shippers with the new value. The key is only shown at generation time and is not returned by 'pipeline get'.`,
			Example:   "ionosctl logging-service pipeline key --location de/txl --pipeline-id ID",
			PreCmdRun: preRunKeyCmd,
			CmdRun:    runKeyCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline to generate a key for", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func preRunKeyCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagLoggingPipelineId)
}

func runKeyCmd(c *core.CommandConfig) error {
	pipelineId, err := c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineId)
	if err != nil {
		return err
	}

	key, _, err := client.Must().LoggingServiceClient.KeyApi.PipelinesKeyPost(
		context.Background(), pipelineId,
	).Body(
		map[string]interface{}{}, // explicit empty body due to 'Error: body is required and must be specified'
	).Execute()
	if err != nil {
		return err
	}

	c.Msg("%s", key.Key)

	return nil
}
