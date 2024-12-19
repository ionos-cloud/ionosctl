package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterRestoreCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongodb",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a Mongo Cluster using a snapshot",
		LongDesc:  "This command restores a cluster via its snapshot. A cluster can have multiple snapshots. A snapshot is added during the following cases:\nWhen a cluster is created, known as initial sync which usually happens in less than 24 hours.\nAfter a restore.\nA snapshot is a copy of the data in the cluster at a certain time. Every 24 hours, a base snapshot is taken, and every Sunday, a full snapshot is taken. Snapshots are retained for the last seven days; hence, recovery is possible for up to a week from the current date.\nYou can restore from any snapshot as long as it was created with the same or older MongoDB patch version.\nSnapshots are stored in an IONOS S3 Object Storage bucket in the same region as your database. Databases in regions where IONOS S3 Object Storage is not available is backed up to eu-central-2.",
		Example:   "ionosctl dbaas mongo cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>",
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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Restoring Cluster %s with snapshot %s",
				clusterId, snapshotId))

			_, err := client.Must().MongoClient.RestoresApi.ClustersRestorePost(context.Background(), clusterId).
				CreateRestoreRequest(
					ionoscloud.CreateRestoreRequest{
						SnapshotId: &snapshotId,
					},
				).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagSnapshotId, "", "", "The unique ID of the snapshot you want to restore.", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagSnapshotId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoSnapshots(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
