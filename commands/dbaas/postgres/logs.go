package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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
	clusterCmd.AddColsFlag(allClusterLogsCols)

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
	list.AddStringFlag(constants.FlagSince, constants.FlagSinceShort, "", "The start time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both start-time and since are set, start-time will be used.")
	list.AddStringFlag(constants.FlagUntil, constants.FlagUntilShort, "", "The end time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both end-time and until are set, end-time will be used.")
	list.AddStringFlag(constants.FlagStartTime, constants.FlagStartTimeShort, "", "The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(constants.FlagEndTime, constants.FlagEndTimeShort, "", "The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(constants.FlagDirection, "", "BACKWARD", "The direction in which to scan through the logs. The logs are returned in order of the direction.")
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BACKWARD", "FORWARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(constants.FlagLimit, constants.FlagLimitShort, 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Minimum: 1. Maximum: 5000")
	list.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return clusterCmd
}

func PreRunClusterLogsList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSince)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, constants.FlagSince)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, constants.FlagSince)), hourSuffix) {
			return errors.New("--since option must have suffix h(hours) or m(minutes). e.g.: --since 2h")
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagUntil)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, constants.FlagUntil)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, constants.FlagUntil)), hourSuffix) {
			return errors.New("--until option must have suffix h(hours) or m(minutes). e.g.: --until 1h")
		}
	}
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterLogsList(c *core.CommandConfig) error {
	c.Verbose("%s: %v", constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))

	queryParams, err := getLogsQueryParams(c)
	if err != nil {
		return err
	}

	c.Verbose("Getting Logs for the specified Cluster...")

	req := client.Must().PostgresClient.LogsApi.
		ClusterLogsGet(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if queryParams != nil {
		if !queryParams.StartTime.IsZero() {
			req = req.Start(queryParams.StartTime)
		}
		if !queryParams.EndTime.IsZero() {
			req = req.End(queryParams.EndTime)
		}
		if queryParams.Limit != 0 {
			req = req.Limit(queryParams.Limit)
		}
		if queryParams.Direction != "" {
			req = req.Direction(queryParams.Direction)
		}
	}
	clusterLogs, _, err := req.Execute()
	if err != nil {
		return err
	}

	// Flatten instances -> messages. Postgres logs group messages by instance,
	// concatenating messages and times with newlines per instance row.
	var rows []map[string]any
	for _, instance := range clusterLogs.GetInstances() {
		if instance.GetMessages() == nil {
			continue
		}
		var messages, times, ls string
		for _, msg := range instance.GetMessages() {
			messages = fmt.Sprintf("%v%v\n", messages, msg.GetMessage())
			times = fmt.Sprintf("%v%v\n", times, msg.GetTime())
			ls = fmt.Sprintf("%vMessage: %v Time:%v\n", ls, msg.GetMessage(), msg.GetTime())
		}
		rows = append(rows, map[string]any{
			"Name":    instance.GetName(),
			"Message": messages,
			"Time":    times,
			"Logs":    ls,
		})
	}

	return c.Printer(allClusterLogsCols).Print(rows)
}

type LogsQueryParams struct {
	Direction          string
	Limit              int32
	StartTime, EndTime time.Time
}

func getLogsQueryParams(c *core.CommandConfig) (*LogsQueryParams, error) {
	var (
		startTime, endTime time.Time
		err                error
	)

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSince)) && !viper.IsSet(core.GetFlagName(c.NS, constants.FlagStartTime)) {
		since := viper.GetString(core.GetFlagName(c.NS, constants.FlagSince))

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

		c.Verbose("Since: %v. StartTime [RFC3339 format]: %v", since, startTime)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagUntil)) && !viper.IsSet(core.GetFlagName(c.NS, constants.FlagEndTime)) {
		until := viper.GetString(core.GetFlagName(c.NS, constants.FlagUntil))

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

		c.Verbose("Until: %v. End Time [RFC3339 format]: %v", until, endTime)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStartTime)) {
		c.Verbose("Start Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagStartTime)))

		startTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagStartTime)))
		if err != nil {
			return nil, err
		}
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagEndTime)) {
		c.Verbose("End Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagEndTime)))

		endTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagEndTime)))
		if err != nil {
			return nil, err
		}
	}

	c.Verbose("Direction: %v", strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagDirection))))
	c.Verbose("Limit: %v", viper.GetInt32(core.GetFlagName(c.NS, constants.FlagLimit)))

	return &LogsQueryParams{
		Direction: strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagDirection))),
		Limit:     viper.GetInt32(core.GetFlagName(c.NS, constants.FlagLimit)),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Output Printing

var allClusterLogsCols = []table.Column{
	{Name: "Logs", JSONPath: "Logs", Default: true},
	{Name: "Name", JSONPath: "Name"},
	{Name: "Message", JSONPath: "Message"},
	{Name: "Time", JSONPath: "Time"},
}
