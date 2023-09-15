package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagStartDuration = "start"
	flagStart         = "startDate"
	flagEndDuration   = "end"
	flagEnd           = "endDate"
	flagDirection     = "direction"
	flagLimit         = "limit"
)

func LogsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List (and optionally filter) the logs of your Mongo Cluster. Use --cols message to see the logs messages.",
		Example:   "ionosctl dbaas mongo logs list --cluster-id CLUSTER_ID --start -24h --end -20h --limit 1 --direction FORWARD --cols message",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}

			c.Command.Command.MarkFlagsMutuallyExclusive(flagStart, flagStartDuration)
			c.Command.Command.MarkFlagsMutuallyExclusive(flagEnd, flagEndDuration)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting logs of Cluster %s", clusterId))

			req := client.Must().MongoClient.LogsApi.ClustersLogsGet(context.Background(), clusterId)
			if fn := core.GetFlagName(c.NS, flagStart); viper.IsSet(fn) {
				start, err := time.Parse(time.RFC3339, viper.GetString(fn))
				if err != nil {
					return fmt.Errorf("failed parsing start time as RFC3339: %w", err)
				}
				req = req.Start(start)
			}
			if fn := core.GetFlagName(c.NS, flagStartDuration); viper.IsSet(fn) {
				start := time.Now().Add(viper.GetDuration(fn))
				req = req.Start(start)
			}

			if fn := core.GetFlagName(c.NS, flagEnd); viper.IsSet(fn) {
				end, err := time.Parse(time.RFC3339, viper.GetString(fn))
				if err != nil {
					return fmt.Errorf("failed parsing end time as RFC3339: %w", err)
				}
				req = req.End(end)
			}
			if fn := core.GetFlagName(c.NS, flagEndDuration); viper.IsSet(fn) {
				end := time.Now().Add(viper.GetDuration(fn))
				req = req.End(end)
			}

			if fn := core.GetFlagName(c.NS, flagLimit); viper.IsSet(fn) {
				limit := viper.GetInt32(fn)
				req = req.Limit(limit)
			}

			if fn := core.GetFlagName(c.NS, flagDirection); viper.IsSet(fn) {
				direction := viper.GetString(fn)
				req = req.Direction(direction)
			}

			logs, _, err := req.Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			logsConverted, err := convertLogsToTable(logs.Instances)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(logs, logsConverted,
				printer.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddDurationFlag(flagStartDuration, "", 0*time.Second, "The start time, as a duration. This should be negative, i.e. -720h. Valid: h, m, s")
	cmd.AddStringFlag(flagStart, "", "", "The start time for the query in RFC3339 format. Must not be greater than 30 days ago and less than the end parameter. The default value is 30 days ago.")
	cmd.AddDurationFlag(flagEndDuration, "", 0*time.Second, "The end time, as a duration. This should be negative and greater than the start time, i.e. -24h. Valid: h, m, s")
	cmd.AddStringFlag(flagEnd, "", "", "The end time for the query in RFC3339 format. Must not be greater than the start parameter. The default value is the current timestamp.")
	cmd.AddSetFlag(flagDirection, "", "", []string{"BACKWARD", "FORWARD"}, "The direction in which to scan through the logs. The logs are returned in order of the direction")
	cmd.AddIntFlag(flagLimit, "", 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Must be between 1 - 5000")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
