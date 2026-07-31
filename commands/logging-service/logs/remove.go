package logs

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func LogsRemoveCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "remove",
			ShortDesc: "Remove a log stream from a logging pipeline",
			LongDesc: `Remove one log stream, selected by its --log-tag, and patch the remaining logs back into the pipeline.

NOTE: a pipeline must always contain at least one log, so removing the last remaining log is not allowed.`,
			Example:   `ionosctl logging-service logs remove --location de/txl --pipeline-id ID --log-tag TAG`,
			PreCmdRun: preRunRemoveCmd,
			CmdRun:    runRemoveCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the pipeline containing the log stream", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Tag of the log stream to remove (identifies which log within the pipeline)",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func runRemoveCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete %s", tag), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	var newLogs []logging.PipelineNoAddrLogs
	for _, log := range patchPipeline.Properties.Logs {
		if log.Tag == tag {
			continue
		}

		newLogs = append(newLogs, log)
	}
	patchPipeline.Properties.Logs = newLogs

	_, _, err = client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).PipelinePatch(
		*patchPipeline,
	).Execute()
	if err != nil {
		return err
	}

	c.Msg("Log successfully removed from pipeline")

	return nil
}

func preRunRemoveCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(
		constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}
