package logs

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func LogsGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "get",
			ShortDesc: "Retrieve a single log stream from a logging pipeline",
			LongDesc:  `Show one log stream of a pipeline, selected by its tag: its source, protocol, destination (type and retention) and labels.`,
			Example:   `ionosctl logging-service logs get --location de/txl --pipeline-id ID --log-tag TAG`,
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the pipeline containing the log stream", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Tag of the log stream to retrieve (identifies which log within the pipeline)",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	return cmd
}

func runGetCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	var log logging.PipelineNoAddrLogs

	for _, l := range pipeline.Properties.Logs {
		if l.Tag == tag {
			log = l

			break
		}
	}

	return c.Printer(allCols).Print(log)
}

func preRunGetCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(
		constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}
