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

func LogsUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "update",
			ShortDesc: "Update a log stream in a logging pipeline",
			LongDesc:  `Update one log stream, selected by its --log-tag, and patch it back into the pipeline. Only the flags you pass are changed; every other attribute of the log is preserved. You can rename the tag (--new-log-tag), change the --log-source or --log-protocol, and adjust the destination backend (--log-type) and retention (--log-retention-time).`,
			Example: `# Extend a log stream's retention to 30 days
ionosctl logging-service logs update --location de/txl --pipeline-id ID --log-tag k8s --log-retention-time 30`,
			PreCmdRun: preRunUpdateCmd,
			CmdRun:    runUpdateCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline containing the log stream", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Tag of the log stream to update (identifies which log within the pipeline)",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LoggingServiceLogTags(
				viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId)),
			)
		}, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		"new-"+constants.FlagLoggingPipelineLogTag, "", "", "Rename the log stream's tag. Leave unset to keep the current tag",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogSource, "", "", constants.EnumLogSources,
		"New source for the log stream. One of: docker, systemd, kubernetes, generic. Leave unset to keep the current source",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogProtocol, "", "", constants.EnumLogProtocols,
		"New protocol for the log stream. One of: http, tcp. Leave unset to keep the current protocol",
	)
	cmd.AddStringSliceFlag(constants.FlagLoggingPipelineLogLabels, "", nil, "Labels for the log stream. Comma-separated, e.g. --log-labels env=prod,team=core")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogType, "", "loki",
		"Destination backend the logs are stored in and queried from. Currently 'loki'",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogRetentionTime, "", "30", constants.EnumLogRetentionTime,
		"How many days logs are kept before deletion. One of: 7, 14, 30",
	)

	return cmd
}

func preRunUpdateCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(
		constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
	)
}

func runUpdateCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))

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

	var newLogs []logging.PipelineNoAddrLogs
	for _, log := range patchPipeline.Properties.Logs {
		if log.Tag == tag {
			newLog = fillOutEmptyFields(log, newLog)

			continue
		}

		newLogs = append(newLogs, log)
	}
	newLogs = append(newLogs, newLog)
	patchPipeline.Properties.Logs = newLogs

	newPipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).PipelinePatch(
		*patchPipeline,
	).Execute()
	if err != nil {
		return err
	}

	return handleLogPrint(newPipeline, c)
}
