package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func SnapshotsListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "snapshot",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List the snapshots (backups) of a Mongo Cluster",
		LongDesc:  "List the snapshots of a MongoDB cluster - the point-in-time backups you can restore from with `cluster restore`. Snapshots are retained for the last 7 days. The Version column is the MongoDB version each snapshot was taken on; you can only restore onto a cluster running that version or newer. Playground clusters have no snapshots.",
		Example:   "ionosctl dbaas mongo snapshot list --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Verbose("Getting snapshots of Cluster %s", clusterId)

			snapshots, _, err := client.Must().MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(snapshots)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster whose snapshots to list", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true

	return cmd
}
