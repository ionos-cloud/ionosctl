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
		ShortDesc: "Restore an In-Memory DB Cluster in place to a point in time",
		LongDesc: `Use this command to trigger an in-place restore of the specified In-Memory DB Cluster from its own snapshots.

The cluster is restored to the most recent snapshot taken at or before ` + "`--recovery-time`" + `. Use ` + "`snapshot get`" + ` on one of the cluster's snapshots to see the available recovery window (earliestRecoveryTargetTime / latestRecoveryTargetTime).

To instead create a NEW cluster from a specific snapshot, use ` + "`cluster create --snapshot-id`" + `.

Required values to run command:

* Cluster Id
* Recovery Time`,
		Example:    "ionosctl dbaas in-memory-db-v2 cluster restore --cluster-id <cluster-id> --recovery-time 2024-01-15T10:00:00Z --password <password>",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, "", "", "An ISO 8601 timestamp (RFC3339, e.g. 2024-01-15T10:00:00Z) to restore the cluster to. The nearest snapshot taken at or before this time is used", core.RequiredFlagOption())
	restoreCmd.AddStringFlag(constants.ArgUser, "", "", "Username for the In-Memory DB user. Defaults to the cluster's current username")
	restoreCmd.AddStringFlag(constants.ArgPassword, "", "", "Password for the In-Memory DB user. Required because the API does not return it on GET requests. Plaintext is hashed (SHA-256) client-side", core.RequiredFlagOption())

	return restoreCmd
}

func PreRunClusterRestore(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.FlagRecoveryTime, constants.ArgPassword)
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	recoveryTimeRaw := viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime))

	c.Verbose(constants.ClusterId, clusterId)

	recoveryTime, err := time.Parse(time.RFC3339, recoveryTimeRaw)
	if err != nil {
		return fmt.Errorf("invalid recovery-time format (expected RFC3339, e.g. 2024-01-15T10:00:00Z): %w", err)
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v to %v", clusterId, recoveryTimeRaw), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	clusterProperties := clusterRead.Properties

	// In-place restore of an existing cluster uses inPlaceRestore (a recovery
	// timestamp), NOT sourceSnapshotId — that variant is only valid when creating
	// a new cluster (see `cluster create --snapshot-id`).
	inPlace := inmemorydb.NewInPlaceRestoreClusterFromSnapshot(recoveryTime)
	c.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", recoveryTime)
	clusterProperties.RestoreFromSnapshot = &inmemorydb.ClusterRestoreFromSnapshot{InPlaceRestoreClusterFromSnapshot: inPlace}

	// The API does not return the password on GET, so the fetched cluster carries
	// a password with an empty algorithm that a PUT would reject. Rebuild
	// credentials from the required --password.
	if err := applyCredentialsFromFlags(c, &clusterProperties); err != nil {
		return err
	}

	c.Verbose("Restoring Cluster in place...")

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
