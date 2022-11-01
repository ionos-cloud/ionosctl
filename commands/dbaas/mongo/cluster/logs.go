package cluster

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	flagStart     = "start"
	flagEnd       = "end"
	flagDirection = "direction"
	flagLimit     = "limit"
)

func ClusterLogsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "logs",
		Aliases:   []string{"lg"},
		ShortDesc: "List the logs of your Mongo Cluster",
		Example:   "ionosctl dbaas mongo cluster logs --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			limit := viper.GetInt32(flagLimit)
			direction := viper.GetString(flagDirection)

			start, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, flagStart)))
			if err != nil {
				return err
			}
			end, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, flagEnd)))
			if err != nil {
				return err
			}

			c.Printer.Verbose("Getting logs of Cluster %s", clusterId)
			logsQueryParams := resources.LogsQueryParams{
				StartTime: start,
				EndTime:   end,
				Limit:     limit,
				Direction: direction,
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
		return completer.MongoClusterIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(flagStart, "", "", "The start time for the query in RFC3339 format. Must not be greater than 30 days ago and less than the end parameter. The default value is 30 days ago.")
	cmd.AddStringFlag(flagEnd, "", "", "The end time for the query in RFC3339 format. Must not be greater than the start parameter. The default value is the current timestamp.")
	cmd.AddSetFlag(flagDirection, "", "", []string{"BACKWARD", "FORWARD"}, "The direction in which to scan through the logs. The logs are returned in order of the direction.")
	cmd.AddIntFlag(flagLimit, "", 100, "The maximal number of log lines to return. If the limit is reached then log lines will be cut at the end (respecting the scan direction). Must be between 1 - 5000")

	cmd.Command.SilenceUsage = true

	return cmd
}

type LogsPrint struct {
	InstanceNumber int       `json:"InstanceNumber,omitempty"`
	Name           string    `json:"SnapshotId,omitempty"`
	MessageNumber  int       `json:"MessageNumber,omitempty"`
	Message        string    `json:"Message,omitempty"`
	Time           time.Time `json:"Time,omitempty"`
}

func MakeLogsPrintObject(logs *[]ionoscloud.ClusterLogsInstances) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*logs))
	for idx, instance := range *logs {
		for msgIdx, msg := range *instance.GetMessages() {
			var logsPrint LogsPrint

			logsPrint.InstanceNumber = idx
			logsPrint.MessageNumber = msgIdx
			logsPrint.Name = *instance.GetName()
			logsPrint.Message = *msg.GetMessage()
			logsPrint.Time = *msg.GetTime()

			o := structs.Map(logsPrint)
			out = append(out, o)
		}
	}

	return out
}

func getLogsPrint(c *core.CommandConfig, dcs *[]ionoscloud.ClusterLogsInstances) printer.Result {
	r := printer.Result{}
	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeLogsPrintObject(dcs)                                                                                                 // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(structs.Names(LogsPrint{}), viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
