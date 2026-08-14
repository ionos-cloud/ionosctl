package logs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allCols = []table.Column{
	{Name: "Tag", JSONPath: "tag", Default: true},
	{Name: "Source", JSONPath: "source", Default: true},
	{Name: "Protocol", JSONPath: "protocol", Default: true},
	{Name: "Public", JSONPath: "public", Default: true},
	{Name: "Destinations", Default: true, Format: func(item map[string]any) any {
		dests, ok := item["destinations"].([]any)
		if !ok {
			return nil
		}
		var parts []string
		for _, d := range dests {
			dm, _ := d.(map[string]any)
			if dm == nil {
				continue
			}
			parts = append(parts, fmt.Sprintf("%v (%v days)", dm["type"], dm["retentionInDays"]))
		}
		return strings.Join(parts, ", ")
	}},
	{Name: "Labels", JSONPath: "labels"},
	{Name: "PipelineId", Format: func(item map[string]any) any {
		return item["_pipelineId"]
	}},
}

func LogsCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use: "logs",
			Short: "The subcommands of `ionosctl logging-service logs` allow you to manage logging pipelines. " +
				"They are the backbone of a centralized logging system, " +
				"referring to an instance or configuration of the logging service you can create",
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(LogsListCmd())
	cmd.AddCommand(LogsAddCmd())
	cmd.AddCommand(LogsUpdateCmd())
	cmd.AddCommand(LogsRemoveCmd())
	cmd.AddCommand(LogsGetCmd())
	return cmd
}

// flattenPipelineLogs flattens the logs of all pipelines into a single
// {"items": [...]} payload, tagging each log with its parent pipeline ID so the
// PipelineId column can be rendered. The {"items": ...} shape lets
// [core.CommandConfig.ListAllLocations] merge logs across locations and render
// them (text, json, api-json) uniformly.
func flattenPipelineLogs(pipelines logging.PipelineReadList) map[string]any {
	items := make([]any, 0)
	for _, p := range pipelines.Items {
		for _, log := range p.Properties.Logs {
			b, err := json.Marshal(log)
			if err != nil {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				continue
			}
			m["_pipelineId"] = p.Id
			items = append(items, m)
		}
	}
	return map[string]any{"items": items}
}

func handleLogPrint(pipeline logging.PipelineRead, c *core.CommandConfig) error {
	return c.Printer(allCols).Print(pipeline.Properties.Logs)
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

	if viper.IsSet(core.GetFlagName(c.NS, "new-"+constants.FlagLoggingPipelineLogTag)) {
		newTag = viper.GetString(core.GetFlagName(c.NS, "new-"+constants.FlagLoggingPipelineLogTag))

		newLog.Tag = newTag
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)) {
		source = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)))

		newLog.Source = source
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)) {
		protocol = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)))

		newLog.Protocol = protocol
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)) {
		typ = strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)))

		dest.Type = typ
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime)) {
		retentionTime = viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime))

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
