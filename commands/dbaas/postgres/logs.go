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
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultClusterLogsCols, tabheaders.ColsMessage(allClusterLogsCols))
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
	if c.Command.Command.Flags().Changed(constants.FlagSince) {
		since, err := c.Command.Command.Flags().GetString(constants.FlagSince)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(since, minuteSuffix) &&
			!strings.HasSuffix(since, hourSuffix) {
			return errors.New("--since option must have suffix h(hours) or m(minutes). e.g.: --since 2h")
		}
	}
	if c.Command.Command.Flags().Changed(constants.FlagUntil) {
		until, err := c.Command.Command.Flags().GetString(constants.FlagUntil)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(until, minuteSuffix) &&
			!strings.HasSuffix(until, hourSuffix) {
			return errors.New("--until option must have suffix h(hours) or m(minutes). e.g.: --until 1h")
		}
	}
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterLogsList(c *core.CommandConfig) error {
	clusterId, err := c.Command.Command.Flags().GetString(constants.FlagClusterId)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))

	queryParams, err := getLogsQueryParams(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Logs for the specified Cluster..."))

	req := client.Must().PostgresClient.LogsApi.
		ClusterLogsGet(context.Background(), clusterId)
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

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	logsConverted, err := resource2table.ConvertDbaasPostgresLogsToTable(clusterLogs.Instances)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		clusterLogs, logsConverted,
		tabheaders.GetHeaders(allClusterLogsCols, defaultClusterLogsCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
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

	startTimeChanged := c.Command.Command.Flags().Changed(constants.FlagStartTime)
	if c.Command.Command.Flags().Changed(constants.FlagSince) && !startTimeChanged {
		since, err := c.Command.Command.Flags().GetString(constants.FlagSince)
		if err != nil {
			return nil, err
		}

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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Since: %v. StartTime [RFC3339 format]: %v", since, startTime))
	}

	endTimeChanged := c.Command.Command.Flags().Changed(constants.FlagEndTime)
	if c.Command.Command.Flags().Changed(constants.FlagUntil) && !endTimeChanged {
		until, err := c.Command.Command.Flags().GetString(constants.FlagUntil)
		if err != nil {
			return nil, err
		}

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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Until: %v. End Time [RFC3339 format]: %v", until, endTime))
	}

	if startTimeChanged {
		startTimeStr, err := c.Command.Command.Flags().GetString(constants.FlagStartTime)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Start Time [RFC3339 format]: %v", startTimeStr))

		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, err
		}
	}

	if endTimeChanged {
		endTimeStr, err := c.Command.Command.Flags().GetString(constants.FlagEndTime)
		if err != nil {
			return nil, err
		}
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("End Time [RFC3339 format]: %v", endTimeStr))

		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, err
		}
	}

	direction, err := c.Command.Command.Flags().GetString(constants.FlagDirection)
	if err != nil {
		return nil, err
	}
	direction = strings.ToUpper(direction)

	limit, err := c.Command.Command.Flags().GetInt(constants.FlagLimit)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Direction: %v", direction))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Limit: %v", limit))

	return &LogsQueryParams{
		Direction: direction,
		Limit:     int32(limit),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Output Printing

var (
	defaultClusterLogsCols = []string{"Logs"}
	allClusterLogsCols     = []string{"Name", "Message", "Time", "Logs"}
)
