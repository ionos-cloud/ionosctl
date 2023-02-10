package logs

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/services/dbaas-mongo/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
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
		ShortDesc: "List the logs of your Mongo Cluster",
		Example:   "ionosctl dbaas mongo logs list --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}

			c.Command.Command.MarkFlagsMutuallyExclusive()

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			var startPtr *time.Time = nil
			if viper.IsSet(core.GetFlagName(c.NS, flagStart)) {
				start, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, flagStart)))
				if err != nil {
					return err
				}
				startPtr = &start
			}
			if viper.IsSet(core.GetFlagName(c.NS, flagStartDuration)) {
				start := time.Now().Add(viper.GetDuration(core.GetFlagName(c.NS, flagStartDuration)))
				startPtr = &start
			}

			var endPtr *time.Time = nil
			if viper.IsSet(core.GetFlagName(c.NS, flagEnd)) {
				end, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, flagEnd)))
				if err != nil {
					return err
				}
				endPtr = &end
			}
			if viper.IsSet(core.GetFlagName(c.NS, flagEndDuration)) {
				end := time.Now().Add(viper.GetDuration(core.GetFlagName(c.NS, flagEndDuration)))
				endPtr = &end
			}

			var limitPtr *int32 = nil
			if viper.IsSet(core.GetFlagName(c.NS, flagLimit)) {
				limit := viper.GetInt32(core.GetFlagName(c.NS, flagLimit))
				limitPtr = &limit
			}

			var directionPtr *string = nil
			if viper.IsSet(core.GetFlagName(c.NS, flagDirection)) {
				direction := viper.GetString(core.GetFlagName(c.NS, flagDirection))
				directionPtr = &direction
			}

			c.Printer.Verbose("Getting logs of Cluster %s", clusterId)
			logsQueryParams := resources.LogsQueryParams{
				StartTime: startPtr,
				EndTime:   endPtr,
				Limit:     limitPtr,
				Direction: directionPtr,
			}
			logs, _, err := c.DbaasMongoServices.Clusters().LogsList(clusterId, logsQueryParams)
			if err != nil {
				return err
			}
			return c.Printer.Print(getLogsPrint(c, logs.GetInstances()))
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

	cmd.AddDurationFlag(flagStartDuration, "", 0*time.Second /* pflag absolutely insists on being passed a default */, "The start time, as a duration. This should be negative, i.e. -720h. Valid: h, m, s")
	cmd.AddStringFlag(flagStart, "", "", "The start time for the query in RFC3339 format. Must not be greater than 30 days ago and less than the end parameter. The default value is 30 days ago.")
	cmd.AddDurationFlag(flagEndDuration, "", 0*time.Second, "The end time, as a duration. This should be negative and greater than the start time, i.e. -24h. Valid: h, m, s")
	cmd.AddStringFlag(flagEnd, "", "", "The end time for the query in RFC3339 format. Must not be greater than the start parameter. The default value is the current timestamp.")
	cmd.AddSetFlag(flagDirection, "", "", []string{"BACKWARD", "FORWARD"}, "The direction in which to scan through the logs. The logs are returned in order of the direction")
	cmd.AddIntFlag(flagLimit, "", 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Must be between 1 - 5000")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
