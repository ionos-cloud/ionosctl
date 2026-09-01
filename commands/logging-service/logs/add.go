package logs

import (
	"context"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func LogsAddCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "add",
			ShortDesc: "Add a log stream to a logging pipeline",
			LongDesc: `Add a new log stream to an existing pipeline. The log is appended to the pipeline's current logs (nothing existing is replaced), so the tag must be unique within the pipeline.

Valid values:
  --log-source:   docker, systemd, kubernetes, generic
  --log-protocol: http or tcp
  --log-type:     destination backend, currently loki
  --log-retention-time: 7, 14 or 30 days`,
			Example: `# Add a systemd log stream shipped over HTTP
ionosctl logging-service logs add --location de/txl --pipeline-id ID --log-tag sys --log-source systemd --log-protocol http

# Advanced: 14-day retention with labels
ionosctl logging-service logs add --location de/txl --pipeline-id ID --log-tag k8s --log-source kubernetes --log-protocol tcp --log-retention-time 14 --log-labels env=staging`,
			PreCmdRun: preRunAddCmd,
			CmdRun:    runAddCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline to add the log stream to", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Tag identifying this log stream within the pipeline; must be unique and is how you reference the log later", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogSource, "", "", constants.EnumLogSources,
		"The kind of workload producing the logs, which selects how they are parsed. One of: docker, systemd, kubernetes, generic", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogProtocol, "", "", constants.EnumLogProtocols,
		"Transport used to push logs: 'http' targets the pipeline's HTTP address, 'tcp' its TCP address", core.RequiredFlagOption(),
	)
	cmd.AddStringSliceFlag(constants.FlagLoggingPipelineLogLabels, "", nil, "Optional labels attached to every entry of this stream for filtering in Grafana. Comma-separated, e.g. --log-labels env=prod,team=core")
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

func preRunAddCmd(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(
		constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
		constants.FlagLoggingPipelineLogSource, constants.FlagLoggingPipelineLogProtocol,
	)
}

func runAddCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))
	source := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)))
	protocol := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)))
	typ := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)))
	retentionTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime))

	retentionTimeInt, err := strconv.ParseInt(retentionTime, 10, 32)
	if err != nil {
		return err
	}

	retentionTimeInt32 := int32(retentionTimeInt)

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesFindById(
		context.Background(), pId,
	).Execute()
	if err != nil {
		return err
	}

	dest := logging.PipelineNoAddrLogsDestinations{
		Type:            typ,
		RetentionInDays: retentionTimeInt32,
	}

	newLog := logging.PipelineNoAddrLogs{
		Tag:          tag,
		Source:       source,
		Protocol:     protocol,
		Destinations: []logging.PipelineNoAddrLogsDestinations{dest},
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	patchPipeline.Properties.Logs = append(patchPipeline.Properties.Logs, newLog)

	newPipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).PipelinePatch(
		*patchPipeline,
	).Execute()

	return handleLogPrint(newPipeline, c)
}
