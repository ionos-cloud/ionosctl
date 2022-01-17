package postgres

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultClusterLogsCols, printer.ColsMessage(allClusterLogsCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(clusterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterLogsList,
		InitClient: true,
	})
	list.AddStringFlag(dbaaspg.ArgStartTime, dbaaspg.ArgStartTimeShort, "", "The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(dbaaspg.ArgEndTime, dbaaspg.ArgEndTimeShort, "", "The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddIntFlag(dbaaspg.ArgLimit, dbaaspg.ArgLimitShort, 0, "The maximal number of log lines to return. The command will print all logs, if this is not set")
	list.AddStringFlag(dbaaspg.ArgClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(dbaaspg.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return clusterCmd
}

func RunClusterLogsList(c *core.CommandConfig) error {
	var (
		startTime, endTime time.Time
		err                error
	)
	c.Printer.Verbose("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId)))
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)) {
		c.Printer.Verbose("Start Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)))
		startTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStartTime)))
		if err != nil {
			return err
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)) {
		c.Printer.Verbose("End Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)))
		endTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgEndTime)))
		if err != nil {
			return err
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLimit)) {
		c.Printer.Verbose("Limit: %v", viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgLimit)))
	}
	c.Printer.Verbose("Getting Logs for the specified Cluster...")
	clusterLogs, _, err := c.CloudApiDbaasPgsqlServices.Logs().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId)),
		viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgLimit)), startTime, endTime)
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterLogsPrint(c, clusterLogs))
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
			r.Columns = getClusterLogsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
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
