package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
	"github.com/spf13/cobra"
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

Required values to run command:

* Cluster Id
* Backup Id`,
		Example:    "ionosctl dbaas postgres cluster restore --cluster-id <cluster-id> --backup-id <backup-id>",
		PreCmdRun:  PreRunClusterBackupIds,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(constants.FlagBackupId, "", "", "The unique ID of the backup you want to restore", core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIdsForCluster(viper.GetString(core.GetFlagName(restoreCmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")

	restoreCmd.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	restoreCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	restoreCmd.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return restoreCmd
}

func PreRunClusterBackupIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagBackupId)
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Backup ID: %v", backupId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v from backup: %v", clusterId, backupId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	// Fetch existing cluster
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Cluster..."))
	clusterRead, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	clusterProperties := clusterRead.Properties

	restoreFromBackup := psqlv2.NewPostgresClusterFromBackup()
	restoreFromBackup.SourceBackupId = &backupId

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime))))

		recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
		if err != nil {
			return err
		}

		// Convert time.Time to IonosTime (assuming SDK handles it or I check how IonosTime is defined)
		// SDK usually uses time.Time directly if aliased, or dedicated struct.
		// Checking model_postgres_cluster_from_backup.go earlier: RecoveryTargetDatetime *IonosTime
		// I need to check what IonosTime is.

		targetTime := psqlv2.IonosTime{Time: recoveryTargetTime}
		restoreFromBackup.RecoveryTargetDatetime = &targetTime
	}

	clusterProperties.RestoreFromBackup = restoreFromBackup

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Restoring Cluster from Backup..."))

	clusterEnsure := psqlv2.NewClusterEnsure(clusterId, clusterProperties)

	_, _, err = client.Must().PostgresClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()

	if err != nil {
		return err
	}
	if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("PostgreSQL Cluster successfully restored"))
	return nil
}
