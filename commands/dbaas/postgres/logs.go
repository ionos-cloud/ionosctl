package postgres

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	minuteSuffix = "m"
	hourSuffix   = "h"
)

func LogsCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "logs",
			Short:            "PostgreSQL Cluster Logs Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres logs` allow you to get the Logs of a specified PostgreSQL Cluster.",
			TraverseChildren: true,
		},
	}
	globalFlags := clusterCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultClusterLogsCols, printer.ColsMessage(allClusterLogsCols))
	_ = viper.BindPFlag(core.GetFlagName(clusterCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterLogsCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres-cluster",
		Resource:   "logs",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Logs for a PostgreSQL Cluster",
		Example:    listLogsExample,
		LongDesc:   "Use this command to retrieve the Logs of a specified PostgreSQL Cluster. By default, the result will contain all Cluster Logs. You can specify the start time, end time or a limit for sorting Cluster Logs.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterLogsList,
		CmdRun:     RunClusterLogsList,
		InitClient: true,
	})
	list.AddStringFlag(dbaaspg.ArgSince, dbaaspg.ArgSinceShort, "", "The start time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both start-time and since are set, start-time will be used.")
	list.AddStringFlag(dbaaspg.ArgUntil, dbaaspg.ArgUntilShort, "", "The end time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both end-time and until are set, end-time will be used.")
	list.AddStringFlag(dbaaspg.ArgStartTime, dbaaspg.ArgStartTimeShort, "", "The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(dbaaspg.ArgEndTime, dbaaspg.ArgEndTimeShort, "", "The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(dbaaspg.ArgDirection, dbaaspg.ArgDirectionShort, "BACKWARD", "The direction in which to scan through the logs. The logs are returned in order of the direction.")
	_ = list.Command.RegisterFlagCompletionFunc(dbaaspg.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BACKWARD", "FORWARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(dbaaspg.ArgLimit, dbaaspg.ArgLimitShort, 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Minimum: 1. Maximum: 5000")
	list.AddUUIDFlag(constants.FlagClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption(), core.CompletionsOption(completer.ClustersIds(os.Stderr)))
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return clusterCmd
}

func PreRunClusterLogsList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgSince)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgSince)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgSince)), hourSuffix) {
			return errors.New("--since option must have suffix h(hours) or m(minutes). e.g.: --since 2h")
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgUntil)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgUntil)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgUntil)), hourSuffix) {
			return errors.New("--until option must have suffix h(hours) or m(minutes). e.g.: --until 1h")
		}
	}
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterLogsList(c *core.CommandConfig) error {
	c.Printer.Verbose("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	queryParams, err := getLogsQueryParams(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Getting Logs for the specified Cluster...")
	clusterLogs, _, err := c.CloudApiDbaasPgsqlServices.Logs().Get(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), queryParams)
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterLogsPrint(c, clusterLogs))
}

func getLogsQueryParams(c *core.CommandConfig) (*resources.LogsQueryParams, error) {
	var (
		startTime, endTime time.Time
		err                error
	)
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgSince)) && !viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)) {
		since := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgSince))
		if strings.Contains(since, hourSuffix) {
			noHours, err := strconv.Atoi(strings.TrimSuffix(since, hourSuffix))
			if err != nil {
				return nil, err
			}
			startTime = time.Now().UTC()
			startTime = startTime.Add(-time.Hour * time.Duration(noHours))
		}
		if strings.Contains(since, minuteSuffix) {
			noMinutes, err := strconv.Atoi(strings.TrimSuffix(since, minuteSuffix))
			if err != nil {
				return nil, err
			}
			startTime = time.Now().UTC()
			startTime = startTime.Add(-time.Minute * time.Duration(noMinutes))
		}
		c.Printer.Verbose("Since: %v. StartTime [RFC3339 format]: %v", since, startTime)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgUntil)) && !viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)) {
		until := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgUntil))
		if strings.Contains(until, hourSuffix) {
			noHours, err := strconv.Atoi(strings.TrimSuffix(until, hourSuffix))
			if err != nil {
				return nil, err
			}
			endTime = time.Now().UTC()
			endTime = endTime.Add(-time.Hour * time.Duration(noHours))
		}
		if strings.Contains(until, minuteSuffix) {
			noMinutes, err := strconv.Atoi(strings.TrimSuffix(until, minuteSuffix))
			if err != nil {
				return nil, err
			}
			endTime = time.Now().UTC()
			endTime = endTime.Add(-time.Minute * time.Duration(noMinutes))
		}
		c.Printer.Verbose("Until: %v. End Time [RFC3339 format]: %v", until, endTime)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)) {
		c.Printer.Verbose("Start Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)))
		startTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)))
		if err != nil {
			return nil, err
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)) {
		c.Printer.Verbose("End Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)))
		endTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)))
		if err != nil {
			return nil, err
		}
	}
	c.Printer.Verbose("Direction: %v", strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDirection))))
	c.Printer.Verbose("Limit: %v", viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgLimit)))
	return &resources.LogsQueryParams{
		Direction: strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDirection))),
		Limit:     viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgLimit)),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Output Printing

var (
	defaultClusterLogsCols = []string{"Logs"}
	allClusterLogsCols     = []string{"Name", "Message", "Time", "Logs"}
)

type ClusterLogsPrint struct {
	Name    string `json:"Name,omitempty"`
	Logs    string `json:"Logs,omitempty"`
	Message string `json:"Message,omitempty"`
	Time    string `json:"Time,omitempty"`
}

func getClusterLogsPrint(c *core.CommandConfig, logs *resources.ClusterLogs) printer.Result {
	r := printer.Result{}
	if c != nil {
		if logs != nil {
			r.OutputJSON = logs
			r.KeyValue = getClusterLogsKVMaps(*logs)
			r.Columns = getClusterLogsCols(core.GetFlagName(c.Resource, constants.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getClusterLogsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultClusterLogsCols
	}
	columnsMap := map[string]string{
		"Name":    "Name",
		"Message": "Message",
		"Time":    "Time",
		"Logs":    "Logs",
	}
	var clusterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			clusterCols = append(clusterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return clusterCols
}

func getClusterLogsKVMaps(clusterLogs resources.ClusterLogs) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)
	if instances, ok := clusterLogs.GetInstancesOk(); ok && instances != nil {
		for _, instance := range *instances {
			var clusterPrint ClusterLogsPrint
			if nameOk, ok := instance.GetNameOk(); ok && nameOk != nil {
				clusterPrint.Name = *nameOk
			}
			if messagesOk, ok := instance.GetMessagesOk(); ok && messagesOk != nil {
				var messages, times string
				var logs string
				for _, msg := range *messagesOk {
					if messageOk, ok := msg.GetMessageOk(); ok && messageOk != nil {
						messages = fmt.Sprintf("%s%s\n", messages, *messageOk)
						logs = fmt.Sprintf("%sMessage: %s ", logs, *messageOk)
					}
					if timeOk, ok := msg.GetTimeOk(); ok && timeOk != nil {
						timeOkRFC := timeOk.Format(time.RFC3339)
						times = fmt.Sprintf("%s%s\n", times, timeOkRFC)
						logs = fmt.Sprintf("%sTime: %s\n", logs, timeOkRFC)
					}
				}
				clusterPrint.Logs = strings.TrimSuffix(logs, "\n")
				clusterPrint.Message = strings.TrimSuffix(messages, "\n")
				clusterPrint.Time = strings.TrimSuffix(times, "\n")
			}
			o := structs.Map(clusterPrint)
			out = append(out, o)
		}
	}
	return out
}
