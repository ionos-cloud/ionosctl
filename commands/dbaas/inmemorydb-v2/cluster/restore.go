package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/spf13/viper"
)

func ClusterRestoreCmd() *core.Command {
	ctx := context.TODO()
	restoreCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-inmemorydb-v2",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore an In-Memory DB Cluster from a snapshot",
		LongDesc: `Use this command to restore the specified In-Memory DB Cluster from a snapshot.

Required values to run command:

* Cluster Id
* Snapshot Id`,
		Example:    "ionosctl dbaas in-memory-db-v2 cluster restore --cluster-id <cluster-id> --snapshot-id <snapshot-id>",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagSnapshotId, "", "", "The unique ID of the snapshot you want to restore from", core.RequiredFlagOption(),
		core.WithCompletion(completer.SnapshotIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, "", "", "An ISO 8601 timestamp to restore from the most recent snapshot taken at or before that time. If empty, the latest available snapshot is used")

	return restoreCmd
}

func PreRunClusterRestore(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.FlagSnapshotId)
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	snapshotId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))

	c.Verbose(constants.ClusterId, clusterId)
	c.Verbose("Snapshot ID: %v", snapshotId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v from snapshot: %v", clusterId, snapshotId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	clusterProperties := clusterRead.Properties

	restore := inmemorydb.NewRestoreClusterFromSnapshot(snapshotId)
	if viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) != "" {
		t, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
		if err != nil {
			return fmt.Errorf("invalid recovery-time format (expected RFC3339, e.g. 2024-01-15T10:00:00Z): %w", err)
		}
		c.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", t)
		restore.RecoveryTargetDatetime = &inmemorydb.IonosTime{Time: t}
	}
	clusterProperties.RestoreFromSnapshot = &inmemorydb.ClusterRestoreFromSnapshot{RestoreClusterFromSnapshot: restore}

	c.Verbose("Restoring Cluster from Snapshot...")

	clusterEnsure := inmemorydb.NewClusterEnsure(clusterId, clusterProperties)

	_, _, err = client.Must().InMemoryDBClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	c.Msg("In-Memory DB Cluster successfully restored")
	return nil
}
