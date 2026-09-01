package logs

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
		LongDesc: `Fetch MongoDB server log lines for a cluster, flattened to one row per message (instance, name, message number, message, time).

Bound the query with a time range given EITHER as absolute RFC3339 timestamps (--startDate/--endDate) OR as relative negative durations from now (--start/--end); the absolute and relative forms of each bound are mutually exclusive. The window may reach at most 30 days into the past. --direction sets scan order (BACKWARD from newest, or FORWARD from oldest) and --limit caps the number of returned lines (1-5000). The message text is hidden by default - add --cols message (or --cols Message) to see it.`,
		Example: "ionosctl dbaas mongo logs list --cluster-id CLUSTER_ID --start -24h --end -20h --limit 1 --direction FORWARD --cols message",
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
			c.Verbose("Getting logs of Cluster %s", clusterId)

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

			if fn := core.GetFlagName(c.NS, flagDirection); viper.IsSet(fn) {
				direction := viper.GetString(fn)
				req = req.Direction(direction)
			}

			logs, _, err := req.Execute()
			if err != nil {
				return err
			}

			// Flatten instances -> messages into individual rows
			var rows []map[string]any
			for idx, instance := range logs.GetInstances() {
				for msgIdx, msg := range instance.GetMessages() {
					rows = append(rows, map[string]any{
						"Instance":      idx,
						"Name":          instance.GetName(),
						"MessageNumber": msgIdx,
						"Message":       msg.GetMessage(),
						"Time":          msg.GetTime(),
					})
				}
			}
			return c.Printer(allCols).Print(rows)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster whose logs to fetch", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddDurationFlag(flagStartDuration, "", 0*time.Second, "Relative start of the window as a negative duration from now, e.g. -720h. Units: h, m, s. Mutually exclusive with --startDate; window may not start more than 30 days ago")
	cmd.AddStringFlag(flagStart, "", "", "Absolute start of the window (RFC3339), e.g. 2024-01-15T10:00:00Z. Must be within the last 30 days and before the end. Mutually exclusive with --start. Defaults to 30 days ago")
	cmd.AddDurationFlag(flagEndDuration, "", 0*time.Second, "Relative end of the window as a negative duration from now, e.g. -24h (must be later than the start). Units: h, m, s. Mutually exclusive with --endDate")
	cmd.AddStringFlag(flagEnd, "", "", "Absolute end of the window (RFC3339). Must be after the start. Mutually exclusive with --end. Defaults to now")
	cmd.AddSetFlag(flagDirection, "", "", []string{"BACKWARD", "FORWARD"}, "Scan order: BACKWARD returns newest-first, FORWARD oldest-first. Determines which end is truncated when --limit is hit")
	cmd.AddIntFlag(flagLimit, "", 100, "Maximum number of log lines to return (1-5000). When reached, lines are cut at the end according to --direction")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
