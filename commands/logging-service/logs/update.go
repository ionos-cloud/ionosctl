package logs

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

func LogsUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "update",
			ShortDesc: "Update a log from a logging pipeline",
			Example: `ionosctl logging-service logs update --pipeline-id ID --log-tag TAG --log-source SOURCE --log
-protocol PROTOCOL`,
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)
	cmd.AddStringFlag(
		constants.FlagLogTag, "", "", "The tag of the pipeline log that you want to update",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		"new-"+constants.FlagLogTag, "", "", "The new tag for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLogSource, "", "", constants.EnumLogSources,
		"Sets the source for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLogProtocol, "", "", constants.EnumLogProtocols,
		"Sets the protocol for the pipeline log",
	)
	cmd.AddStringSliceFlag(constants.FlagLogLabels, "", nil, "Sets the labels for the pipeline log")
	cmd.AddStringFlag(
		constants.FlagLogType, "", "loki",
		"Sets the destination type for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLogRetentionTime, "", "30", constants.EnumLogRetentionTime,
		"Sets the retention time in days for the pipeline log",
	)

	return cmd
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagPipelineId, constants.FlagLogTag,
	)
}

func runUpdateCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLogTag))

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	newLog, err := generatePatchObject(c)
	if err != nil {
		return err
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	var newLogs []logging.PipelineCreatePropertiesLogs
	for _, log := range patchPipeline.Properties.Logs {
		if *log.Tag == tag {
			newLog = fillOutEmptyFields(&log, newLog)

			continue
		}

		newLogs = append(newLogs, log)
	}
	newLogs = append(newLogs, *newLog)
	patchPipeline.Properties.Logs = newLogs

	newPipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).Pipeline(
		*patchPipeline,
	).Execute()
	if err != nil {
		return err
	}

	return handleLogPrint(newPipeline, c)
}
