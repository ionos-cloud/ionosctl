package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
	globalFlags.StringSliceP(constants.FlagCols, "", defaultClusterLogsCols, tabheaders.ColsMessage(allClusterLogsCols))
	_ = viper.BindPFlag(core.GetFlagName(clusterCmd.Name(), constants.FlagCols), globalFlags.Lookup(constants.FlagCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddStringFlag(dbaaspg.FlagSince, dbaaspg.FlagSinceShort, "", "The start time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both start-time and since are set, start-time will be used.")
	list.AddStringFlag(dbaaspg.FlagUntil, dbaaspg.FlagUntilShort, "", "The end time for the query using a time delta since the current moment: 2h - 2 hours ago, 20m - 20 minutes ago. Only hours and minutes are supported, and not at the same time. If both end-time and until are set, end-time will be used.")
	list.AddStringFlag(dbaaspg.FlagStartTime, dbaaspg.FlagStartTimeShort, "", "The start time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(dbaaspg.FlagEndTime, dbaaspg.FlagEndTimeShort, "", "The end time for the query in RFC3339 format. Example: 2021-10-05T11:30:17.45Z")
	list.AddStringFlag(dbaaspg.FlagDirection, dbaaspg.FlagDirectionShort, "BACKWARD", "The direction in which to scan through the logs. The logs are returned in order of the direction.")
	_ = list.Command.RegisterFlagCompletionFunc(dbaaspg.FlagDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BACKWARD", "FORWARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(dbaaspg.FlagLimit, dbaaspg.FlagLimitShort, 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Minimum: 1. Maximum: 5000")
	list.AddUUIDFlag(constants.FlagClusterId, dbaaspg.FlagIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return clusterCmd
}

func PreRunClusterLogsList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagSince)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagSince)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagSince)), hourSuffix) {
			return errors.New("--since option must have suffix h(hours) or m(minutes). e.g.: --since 2h")
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagUntil)) {
		if !strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagUntil)), minuteSuffix) &&
			!strings.HasSuffix(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagUntil)), hourSuffix) {
			return errors.New("--until option must have suffix h(hours) or m(minutes). e.g.: --until 1h")
		}
	}
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterLogsList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	queryParams, err := getLogsQueryParams(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Logs for the specified Cluster..."))

	clusterLogs, _, err := c.CloudApiDbaasPgsqlServices.Logs().Get(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), queryParams)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.FlagCols))

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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func getLogsQueryParams(c *core.CommandConfig) (*resources.LogsQueryParams, error) {
	var (
		startTime, endTime time.Time
		err                error
	)

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagSince)) && !viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagStartTime)) {
		since := viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagSince))

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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Since: %v. StartTime [RFC3339 format]: %v", since, startTime))
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagUntil)) && !viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagEndTime)) {
		until := viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagUntil))

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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Until: %v. End Time [RFC3339 format]: %v", until, endTime))
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagStartTime)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Start Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagStartTime))))

		startTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagStartTime)))
		if err != nil {
			return nil, err
		}
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.FlagEndTime)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("End Time [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagEndTime))))

		endTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagEndTime)))
		if err != nil {
			return nil, err
		}
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Direction: %v", strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagDirection)))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Limit: %v", viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.FlagLimit))))

	return &resources.LogsQueryParams{
		Direction: strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, dbaaspg.FlagDirection))),
		Limit:     viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.FlagLimit)),
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

// Output Printing

var (
	defaultClusterLogsCols = []string{"Logs"}
	allClusterLogsCols     = []string{"Name", "Message", "Time", "Logs"}
)
