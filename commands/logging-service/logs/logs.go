package logs

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultCols = []string{"Tag", "Source", "Protocol", "Public", "Destinations"}
	allCols     = []string{"Tag", "Source", "Protocol", "Public", "Destinations", "Labels", "PipelineId"}
)

func LogsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use: "logs",
			Short: "The subcommands of `ionosctl logging-service logs` allow you to manage logging pipelines. " +
				"They are the backbone of a centralized logging system, " +
				"referring to an instance or configuration of the logging service you can create",
		},
	}

	cmd.AddCommand(LogsListCmd())
	cmd.AddCommand(LogsAddCmd())
	cmd.AddCommand(LogsUpdateCmd())
	cmd.AddCommand(LogsRemoveCmd())
	cmd.AddCommand(LogsGetCmd())
	return cmd
}

func handleLogsPrint(pipelines ionoscloud.PipelineListResponse, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	var logs []ionoscloud.PipelineResponse
	for _, p := range *pipelines.Items {
		logs = append(logs, *p.Properties.Logs...)
	}

	logsConverted, err := resource2table.ConvertLoggingServicePipelinesLogsToTable(pipelines)
	if err != nil {
		return fmt.Errorf("could not convert Logging Service Pipeline Logs to table format: %w", err)
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		logs, logsConverted, tabheaders.GetHeaders(
			allCols,
			defaultCols, cols,
		),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func handleLogPrint(pipeline ionoscloud.Pipeline, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	logsConverted, err := resource2table.ConvertLoggingServicePipelineLogsToTable(pipeline)
	if err != nil {
		return fmt.Errorf("could not convert Logging Service Pipeline Logs to table format: %w", err)
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		pipeline.Properties.Logs, logsConverted, tabheaders.GetHeaders(
			allCols,
			defaultCols, cols,
		),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func convertResponsePipelineToPatchRequest(pipeline ionoscloud.Pipeline) (*ionoscloud.PipelinePatch, error) {
	properties, ok := pipeline.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline properties")
	}

	logs, ok := properties.GetLogsOk()
	if !ok || logs == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline Logs")
	}

	var newLogs []ionoscloud.PipelineCreatePropertiesLogs
	for _, log := range *logs {
		l := ionoscloud.PipelineCreatePropertiesLogs{
			Tag:          log.Tag,
			Source:       log.Source,
			Protocol:     log.Protocol,
			Labels:       log.Labels,
			Destinations: log.Destinations,
		}

		newLogs = append(newLogs, l)
	}

	patch := ionoscloud.PipelinePatch{
		Properties: &ionoscloud.PipelinePatchProperties{
			Name: properties.Name,
			Logs: &newLogs,
		},
	}

	return &patch, nil
}

func generatePatchObject(c *core.CommandConfig) (*ionoscloud.PipelineCreatePropertiesLogs, error) {
	var newTag, source, protocol, typ, retentionTime string
	var labels []string
	var retentionTimeInt32 int32

	dest := ionoscloud.Destination{}
	newLog := ionoscloud.PipelineCreatePropertiesLogs{}

	if viper.IsSet(core.GetFlagName(c.NS, "new-"+constants.FlagLoggingPipelineLogTag)) {
		newTag = viper.GetString(core.GetFlagName(c.NS, "new-"+constants.FlagLoggingPipelineLogTag))

		newLog.Tag = &newTag
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)) {
		source = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)))

		newLog.Source = &source
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)) {
		protocol = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)))

		newLog.Protocol = &protocol
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogLabels)) {
		labels = viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogLabels))

		newLog.Labels = &labels
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)) {
		typ = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)))

		dest.Type = &typ
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime)) {
		retentionTime = viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime))

		retentionTimeInt, err := strconv.ParseInt(retentionTime, 10, 32)
		if err != nil {
			return nil, err
		}

		retentionTimeInt32 = int32(retentionTimeInt)
		dest.RetentionInDays = &retentionTimeInt32
	}

	if dest.Type != nil || dest.RetentionInDays != nil {
		newLog.Destinations = &[]ionoscloud.Destination{dest}
	}

	return &newLog, nil
}

func fillOutEmptyFields(oldLog, newLog *ionoscloud.PipelineCreatePropertiesLogs) *ionoscloud.PipelineCreatePropertiesLogs {
	if newLog.Tag == nil {
		newLog.Tag = oldLog.Tag
	}

	if newLog.Source == nil {
		newLog.Source = oldLog.Source
	}

	if newLog.Protocol == nil {
		newLog.Protocol = oldLog.Protocol
	}

	if newLog.Labels == nil {
		newLog.Labels = oldLog.Labels
	}

	if newLog.Destinations == nil {
		newLog.Destinations = oldLog.Destinations
	} else {
		if (*newLog.Destinations)[0].Type == nil {
			(*newLog.Destinations)[0].Type = (*oldLog.Destinations)[0].Type
		}

		if (*newLog.Destinations)[0].RetentionInDays == nil {
			(*newLog.Destinations)[0].RetentionInDays = (*oldLog.Destinations)[0].RetentionInDays
		}
	}

	return newLog
}
