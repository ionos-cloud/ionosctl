package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterRestoreCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongodb",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a Mongo Cluster in place from one of its snapshots",
		LongDesc: `Roll a MongoDB cluster back in place to the state captured by one of its own snapshots. This overwrites the current data of the cluster identified by --cluster-id with the contents of --snapshot-id (list a cluster's snapshots with ` + "`snapshot list --cluster-id <id>`" + `).

How snapshots accumulate: an initial snapshot is taken when the cluster is created (the initial sync, usually within 24h); another is created after each restore; thereafter a base snapshot is taken every 24h and a full snapshot every Sunday. Snapshots are retained for the last 7 days, so recovery is possible up to a week back. Playground clusters have no snapshots and cannot be restored.

Constraints:
  - You can only restore from a snapshot whose MongoDB version is the same as, or older (by patch) than, the cluster's current version.
  - Snapshots are stored in IONOS S3 Object Storage in the cluster's region (eu-central-2 where S3 is unavailable).

Enterprise clusters additionally support point-in-time recovery WITHIN a snapshot's window via the API's recoveryTargetTime; this CLI restore targets a whole snapshot.`,
		Example: "ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			return c.Command.Command.MarkFlagRequired(constants.FlagSnapshotId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			snapshotId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))

			c.Verbose("Restoring Cluster %s with snapshot %s", clusterId, snapshotId)

			_, err := client.Must().MongoClient.RestoresApi.ClustersRestorePost(context.Background(), clusterId).
				CreateRestoreRequest(
					mongo.CreateRestoreRequest{
						SnapshotId: &snapshotId,
					},
				).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster to restore in place (its current data will be overwritten)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagSnapshotId, "", "", "The snapshot to restore from. Must belong to this cluster and its MongoDB version must be the same or older than the cluster's. List options with `snapshot list --cluster-id <id>`", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagSnapshotId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoSnapshots(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
