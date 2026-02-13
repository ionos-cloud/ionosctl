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
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/cobra"
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

func handleLogsPrint(pipelines logging.PipelineReadList, c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	var logs []logging.PipelineNoAddrLogs
	for _, p := range pipelines.Items {
		logs = append(logs, p.Properties.Logs...)
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func handleLogPrint(pipeline logging.PipelineRead, c *core.CommandConfig) error {
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func convertResponsePipelineToPatchRequest(pipeline logging.PipelineRead) (*logging.PipelinePatch, error) {
	properties, ok := pipeline.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline properties")
	}

	logs, ok := properties.GetLogsOk()
	if !ok || logs == nil {
		return nil, fmt.Errorf("could not retrieve Logging Service Pipeline Logs")
	}

	var newLogs []logging.PipelineNoAddrLogs
	for _, log := range logs {
		l := logging.PipelineNoAddrLogs{
			Tag:          log.Tag,
			Source:       log.Source,
			Protocol:     log.Protocol,
			Destinations: log.Destinations,
		}

		newLogs = append(newLogs, l)
	}

	patch := logging.PipelinePatch{
		Properties: logging.PipelineNoAddr{
			Name: properties.Name,
			Logs: newLogs,
		},
	}

	return &patch, nil
}

func generatePatchObject(c *core.CommandConfig) (logging.PipelineNoAddrLogs, error) {
	var newTag, source, protocol, typ, retentionTime string
	var retentionTimeInt32 int32

	dest := logging.PipelineNoAddrLogsDestinations{}
	newLog := logging.PipelineNoAddrLogs{}

	if c.Command.Command.Flags().Changed("new-" + constants.FlagLoggingPipelineLogTag) {
		newTag, _ = c.Command.Command.Flags().GetString("new-" + constants.FlagLoggingPipelineLogTag)

		newLog.Tag = newTag
	}

	if c.Command.Command.Flags().Changed(constants.FlagLoggingPipelineLogSource) {
		source, _ = c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineLogSource)
		source = strings.ToLower(source)

		newLog.Source = source
	}

	if c.Command.Command.Flags().Changed(constants.FlagLoggingPipelineLogProtocol) {
		protocol, _ = c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineLogProtocol)
		protocol = strings.ToLower(protocol)

		newLog.Protocol = protocol
	}

	if c.Command.Command.Flags().Changed(constants.FlagLoggingPipelineLogType) {
		typ, _ = c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineLogType)
		typ = strings.ToLower(typ)

		dest.Type = typ
	}

	if c.Command.Command.Flags().Changed(constants.FlagLoggingPipelineLogRetentionTime) {
		retentionTime, _ = c.Command.Command.Flags().GetString(constants.FlagLoggingPipelineLogRetentionTime)

		retentionTimeInt, err := strconv.ParseInt(retentionTime, 10, 32)
		if err != nil {
			return logging.PipelineNoAddrLogs{}, err
		}

		retentionTimeInt32 = int32(retentionTimeInt)
		dest.RetentionInDays = retentionTimeInt32
	}

	newLog.Destinations = []logging.PipelineNoAddrLogsDestinations{dest}

	return newLog, nil
}

func fillOutEmptyFields(oldLog, newLog logging.PipelineNoAddrLogs) logging.PipelineNoAddrLogs {
	newLog.Tag = oldLog.Tag
	newLog.Source = oldLog.Source
	newLog.Protocol = oldLog.Protocol

	if newLog.Destinations == nil {
		newLog.Destinations = oldLog.Destinations
	} else {
		newLog.Destinations[0].Type = oldLog.Destinations[0].Type
		newLog.Destinations[0].RetentionInDays = oldLog.Destinations[0].RetentionInDays
	}

	return newLog
}
