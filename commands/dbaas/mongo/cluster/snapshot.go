package cluster

import (
	"context"
	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/viper"
	"os"
	"time"

	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/spf13/cobra"
)

func ClusterSnapshotsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "snapshots",
		Aliases:   []string{"r"},
		ShortDesc: "List the snapshots of your Mongo Cluster",
		Example:   "ionosctl dbaas mongo cluster snapshots --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId, err := c.Command.Command.Flags().GetString(constants.FlagClusterId)
			if err != nil {
				return err
			}
			c.Printer.Verbose("Getting snapshots of Cluster %s", clusterId)
			snapshots, _, err := c.DbaasMongoServices.Clusters().SnapshotsList(clusterId)
			if err != nil {
				return err
			}
			return c.Printer.Print(getSnapshotPrint(c, snapshots.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdP, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}

type SnapshotPrint struct {
	SnapshotId   string    `json:"SnapshotId,omitempty"`
	CreationTime time.Time `json:"CreationTime,omitempty"`
	Size         int32     `json:"Size,omitempty"`
	Version      string    `json:"Version,omitempty"`
}

func MakeSnapshotPrintObject(snapshots *[]ionoscloud.SnapshotResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*snapshots))
	for _, snapshot := range *snapshots {
		var snapshotPrint SnapshotPrint
		snapshotPrint.SnapshotId = *snapshot.GetId()
		snapshotPrint.CreationTime = *snapshot.GetProperties().GetCreationTime()
		snapshotPrint.Size = *snapshot.GetProperties().GetSize()
		snapshotPrint.Version = *snapshot.GetProperties().GetVersion()
	}
	return out
}

func getSnapshotPrint(c *core.CommandConfig, dcs *[]ionoscloud.SnapshotResponse) printer.Result {
	r := printer.Result{}
	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = MakeSnapshotPrintObject(dcs)                                                                                                 // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(structs.Names(SnapshotPrint{}), viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
