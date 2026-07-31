package postgres

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allBackupCols = []table.Column{
	{Name: "BackupId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "Active", JSONPath: "properties.active", Default: true},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate", Default: true},
	{Name: "EarliestRecoveryTargetTime", JSONPath: "properties.earliestRecoveryTargetTime", Default: true},
	{Name: "Version", JSONPath: "properties.version"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func BackupCmd() *core.Command {
	ctx := context.TODO()
	backupCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "backup",
			Aliases: []string{"b"},
			Short:   "List and inspect PostgreSQL backups",
			Long: `Inspect PostgreSQL backups. Backups are created automatically by the DBaaS service for every cluster; you cannot create or delete them manually here.

Each backup covers a continuous recovery WINDOW (not a single instant), which is what enables point-in-time recovery. The window's start is shown as EarliestRecoveryTargetTime. A backup is tied to the PostgreSQL version of the cluster that produced it, which matters when restoring or cloning.

Use these backups with 'cluster restore --backup-id' (restore in place) or 'cluster create --backup-id' (clone into a new cluster).`,
			TraverseChildren: true,
		},
	}
	backupCmd.AddColsFlag(allBackupCols)

	/*
		List Command
	*/
	list := core.NewCommand(ctx, backupCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List all backups",
		LongDesc:   "Retrieve every PostgreSQL backup in the account, across all clusters. To see only the backups of one cluster, use 'dbaas postgres cluster backup list --cluster-id'.",
		Example:    listBackupExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunBackupList,
		InitClient: true,
	})
	_ = list // Actually used - added through "NewCommand" func. TODO: This is confusing!

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, backupCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "backup",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a backup",
		Example:    getBackupExample,
		LongDesc:   "Retrieve details of a single backup by its ID, including the cluster it belongs to, its PostgreSQL version, and the start of its recovery window (EarliestRecoveryTargetTime) which bounds the --recovery-time you can use when restoring.\n\nRequired values to run command:\n\n* Backup Id",
		PreCmdRun:  PreRunBackupId,
		CmdRun:     RunBackupGet,
		InitClient: true,
	})
	get.AddStringFlag(constants.FlagBackupId, constants.FlagIdShort, "", "ID of the backup to retrieve. See 'dbaas postgres backup list'", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return backupCmd
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupId)
}

func RunBackupList(c *core.CommandConfig) error {
	c.Verbose("Getting Backups...")

	backups, _, err := client.Must().PostgresClient.BackupsApi.ClustersBackupsGet(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("could not get Backups: %w", err)
	}

	return c.Printer(allBackupCols).Prefix("items").Print(backups)
}

func RunBackupGet(c *core.CommandConfig) error {
	c.Verbose("Backup ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId)))
	c.Verbose("Getting Backup...")

	backup, _, err := client.Must().PostgresClient.BackupsApi.ClustersBackupsFindById(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allBackupCols).Print(backup)
}

func ClusterBackupCmd() *core.Command {
	ctx := context.TODO()
	clusterBackupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "List the backups of a specific cluster",
			Long:             "List the automated PostgreSQL backups belonging to one specific cluster. This is the per-cluster view; 'dbaas postgres backup list' lists backups across all clusters.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterBackupCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres.cluster",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List the backups of one cluster",
		LongDesc:   "Retrieve the backups belonging to a single PostgreSQL cluster, identified by --cluster-id.\n\nRequired values to run command:\n\n* Cluster Id",
		Example:    listBackupExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterBackupList,
		InitClient: true,
	})
	list.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	return clusterBackupCmd
}

func RunClusterBackupList(c *core.CommandConfig) error {
	c.Verbose("%s: %v", constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	c.Verbose("Getting Backups from Cluster...")

	backups, _, err := client.Must().PostgresClient.BackupsApi.ClusterBackupsGet(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).Execute()
	if err != nil {
		return err
	}
	return c.Printer(allBackupCols).Prefix("items").Print(backups)
}
