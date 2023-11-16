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
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
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
		constants.FlagLoggingPipelineId, constants.FlagIdShort, "",
		"The ID of the logging pipeline", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLoggingPipelineId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.LoggingServicePipelineIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Sets the tag for the pipeline log", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogSource, "", "", []string{"docker", "systemd", "generic", "kubernetes"},
		"Sets the source for the pipeline log", core.RequiredFlagOption(),
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogProtocol, "", "", []string{"http", "tcp"},
		"Sets the protocol for the pipeline log", core.RequiredFlagOption(),
	)
	cmd.AddStringSliceFlag(constants.FlagLoggingPipelineLogLabels, "", nil, "Sets the labels for the pipeline log")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogType, "", "loki",
		"Sets the destination type for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogRetentionTime, "", "30", []string{"7", "14", "30"},
		"Sets the retention time in days for the pipeline log",
	)

	return cmd
}

func preRunAddCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagLoggingPipelineId, constants.FlagLoggingPipelineLogTag,
		constants.FlagLoggingPipelineLogSource, constants.FlagLoggingPipelineLogProtocol,
	)
}

func runAddCmd(c *core.CommandConfig) error {
	pId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineId))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))
	source := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)))
	protocol := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)))
	labels := viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogLabels))
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

	dest := ionoscloud.Destination{
		Type:            &typ,
		RetentionInDays: &retentionTimeInt32,
	}

	newLog := ionoscloud.PipelineCreatePropertiesLogs{
		Tag:          &tag,
		Source:       &source,
		Protocol:     &protocol,
		Labels:       &labels,
		Destinations: &[]ionoscloud.Destination{dest},
	}

	patchPipeline, err := convertResponsePipelineToPatchRequest(pipeline)
	if err != nil {
		return err
	}

	*patchPipeline.Properties.Logs = append(*patchPipeline.Properties.Logs, newLog)

	newPipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPatch(
		context.Background(),
		pId,
	).Pipeline(
		*patchPipeline,
	).Execute()

	return handleLogPrint(newPipeline, c)
}
