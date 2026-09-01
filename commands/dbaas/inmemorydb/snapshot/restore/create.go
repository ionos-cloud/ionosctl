package restore

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "restore",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Restore a Replica Set from a Snapshot",
		LongDesc: `Restore an existing In-Memory DB Replica Set from one of its point-in-time snapshots.

You pick the source snapshot with --snapshot-id and the replica set to restore it onto with --replicaset-id; that replica set's data is rolled back to the snapshot's state. The snapshot and the replica set must live in the same location/datacenter. Optionally attach a --name and --description to label the restore operation.`,
		Example: `# Restore a replica set from one of its snapshots
ionosctl dbaas in-memory-db snapshot restore create --snapshot-id SNAPSHOT_ID --replicaset-id REPLICASET_ID

# Restore and label the operation
ionosctl dbaas in-memory-db snapshot restore create --snapshot-id SNAPSHOT_ID --replicaset-id REPLICASET_ID --name nightly-rollback --description "roll back after bad deploy"`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagReplicasetID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := inmemorydb.Restore{}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagReplicasetID)) {
				input.ReplicasetId = viper.GetString(core.GetFlagName(c.NS, constants.FlagReplicasetID))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
				input.DisplayName = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
				input.Description = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription)))
			}

			restore, _, err := client.Must().InMemoryDBClient.RestoreApi.SnapshotsRestoresPost(
				context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))).
				RestoreCreate(inmemorydb.RestoreCreate{Properties: input}).Execute()

			if err != nil {
				return fmt.Errorf("error creating restore: %w", err)
			}

			return c.Printer(allCols).Print(restore)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagSnapshotId, "", "",
		"The ID of the source Snapshot to restore from", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				// for each snapshot
				return utils.SnapshotProperty(func(snapshot inmemorydb.SnapshotRead) string {
					// return its ID
					return snapshot.Id + "\t" + snapshot.Metadata.SnapshotTime.Format("2006-01-02 15:04:05")
				})
			}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations,
		),
		core.RequiredFlagOption(),
	)

	cmd.AddStringFlag(constants.FlagReplicasetID, "", "",
		"The ID of the target Replica Set to restore the snapshot onto (must be in the same location as the snapshot)",
		core.WithCompletion(utils.ReplicasetIDs, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
		core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Optional human-readable name to label this restore operation")
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Optional free-text description of this restore operation")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
