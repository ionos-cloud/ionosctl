package logs

import (
	"context"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

func LogsAddCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "logs",
			Verb:      "add",
			ShortDesc: "Add a log to a logging pipeline",
			Example: `ionosctl logging-service logs add --pipeline-id ID --log-tag TAG --log-source SOURCE --log
-protocol PROTOCOL`,
			PreCmdRun: preRunAddCmd,
			CmdRun:    runAddCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
		core.WithCompletion(completer.LoggingServicePipelineIds, constants.LoggingApiRegionalURL, constants.LoggingLocations),
	)

	cmd.AddStringFlag(
		constants.FlagLogTag, "", "", "Sets the tag for the pipeline log", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLogSource, "", "", constants.EnumLogSources,
		"Sets the source for the pipeline log", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLogProtocol, "", "", constants.EnumLogProtocols,
		"Sets the protocol for the pipeline log", core.RequiredFlagOption(),
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

func preRunAddCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagPipelineId, constants.FlagLogTag,
		constants.FlagLogSource, constants.FlagLogProtocol,
	)
}

func runAddCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLogTag))
	source := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLogSource)))
	protocol := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLogProtocol)))
	labels := viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagLogLabels))
	typ := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLogType)))
	retentionTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagLogRetentionTime))

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

	dest := logging.Destination{
		Type:            &typ,
		RetentionInDays: &retentionTimeInt32,
	}

	newLog := logging.PipelineCreatePropertiesLogs{
		Tag:          &tag,
		Source:       &source,
		Protocol:     &protocol,
		Labels:       labels,
		Destinations: []logging.Destination{dest},
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	patchPipeline.Properties.Logs = append(patchPipeline.Properties.Logs, newLog)

	newPipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).Pipeline(
		*patchPipeline,
	).Execute()

	return handleLogPrint(newPipeline, c)
}
