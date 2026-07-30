package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/viper"
)

func ClusterRestoreCmd() *core.Command {
	ctx := context.TODO()
	restoreCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a PostgreSQL Cluster",
		LongDesc: `Use this command to trigger an in-place restore of the specified PostgreSQL Cluster.

The restore uses the cluster's existing backups automatically; no source backup can be specified. The current data is overwritten with the restored data, and the cluster may experience a brief period of downtime during this process.

Required values to run command:

* Cluster Id
* Recovery Time
* DB Password`,
		Example:    "ionosctl dbaas postgres-v2 cluster restore --cluster-id <cluster-id> --recovery-time <RFC3339-timestamp> --db-password <password>",
		PreCmdRun:  PreRunClusterRestore,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	restoreCmd.AddStringFlag(constants.FlagDbPassword, constants.FlagDbPasswordShortPsql, "", "Password for the initial postgres user", core.RequiredFlagOption())
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "ISO 8601 timestamp up to which the cluster is restored using its existing backups", core.RequiredFlagOption())

	return restoreCmd
}

func PreRunClusterRestore(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.FlagRecoveryTime, constants.FlagDbPassword)
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	c.Verbose(constants.ClusterId, clusterId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("trigger an in-place restore of cluster with id: %v from its existing backups", clusterId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
	if err != nil {
		return fmt.Errorf("invalid recovery-time format (expected RFC3339, e.g. 2024-01-15T10:00:00Z): %w", err)
	}
	c.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", recoveryTargetTime)

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	clusterProperties := clusterRead.Properties

	// In-place restore (cluster modification) infers the source backup automatically;
	// the recovery target time is required and no source backup id may be specified.
	restoreFromBackup := psqlv2.NewPostgresInPlaceRestoreClusterFromBackup(recoveryTargetTime)

	// Password is required because the API does not return it on GET
	password := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbPassword))
	credentials := clusterProperties.GetCredentials()
	credentials.SetPassword(password)
	clusterProperties.SetCredentials(credentials)

	restore := psqlv2.PostgresInPlaceRestoreClusterFromBackupAsClusterRestoreFromBackup(restoreFromBackup)
	clusterProperties.RestoreFromBackup = &restore

	c.Verbose("Restoring Cluster from Backup...")

	clusterEnsure := psqlv2.NewClusterEnsure(clusterId, clusterProperties)

	_, _, err = client.Must().PostgresClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()

	if err != nil {
		return err
	}

	c.Msg("PostgreSQL Cluster successfully restored")
	return nil
}
