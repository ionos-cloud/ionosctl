package pipeline

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
)

func PipelineKeyCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "key",
			ShortDesc: "Generate a new key for a logging pipeline, " +
				"invalidating the old one. The key is used for authentication when sending logs.",
			Example:   "ionosctl logging-service pipeline key --pipeline-id ID",
			PreCmdRun: preRunKeyCmd,
			CmdRun:    runKeyCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline you want to generate a key for", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func preRunKeyCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagLoggingPipelineId)
}

func runKeyCmd(c *core.CommandConfig) error {
	pipelineId, err := c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineId)
	if err != nil {
		return err
	}

	key, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesKeyPost(
		context.Background(), pipelineId,
	).Execute()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateRawOutput(*key.Key))

	return nil
}
