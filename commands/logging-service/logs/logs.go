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
			Use:   "logs",
			Short: "Manage the log streams inside a logging pipeline",
			Long: `A log is one stream of a pipeline. It binds a source (docker, systemd, kubernetes, generic) and a tag to a shipping protocol (http or tcp) and a destination (Loki, with a retention period), plus optional labels. Logs are identified within a pipeline by their tag.

These commands add, list, get, update and remove logs on an existing pipeline (found via --pipeline-id, scoped to --location). Every change is applied by patching the parent pipeline. A pipeline must keep at least one log, so the last remaining log cannot be removed.`,
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

// fillOutEmptyFields keeps the values the user set on newLog and falls back to
// oldLog only for the fields left unset, so that 'logs update' patches just the
// changed attributes instead of resetting the log to its previous state.
func fillOutEmptyFields(oldLog, newLog logging.PipelineNoAddrLogs) logging.PipelineNoAddrLogs {
	if newLog.Tag == "" {
		newLog.Tag = oldLog.Tag
	}
	if newLog.Source == "" {
		newLog.Source = oldLog.Source
	}
	if newLog.Protocol == "" {
		newLog.Protocol = oldLog.Protocol
	}

	if len(newLog.Destinations) == 0 {
		newLog.Destinations = oldLog.Destinations
	} else if len(oldLog.Destinations) > 0 {
		if newLog.Destinations[0].Type == "" {
			newLog.Destinations[0].Type = oldLog.Destinations[0].Type
		}
		if newLog.Destinations[0].RetentionInDays == 0 {
			newLog.Destinations[0].RetentionInDays = oldLog.Destinations[0].RetentionInDays
		}
	}

	return newLog
}
