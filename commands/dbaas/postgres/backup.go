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
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "PostgreSQL Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres backup` allow you to list, get PostgreSQL Backups.",
			TraverseChildren: true,
		},
	}
	globalFlags := backupCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", nil, table.ColsMessage(allBackupCols))
	_ = viper.BindPFlag(core.GetFlagName(backupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = backupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allBackupCols), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, backupCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Cluster Backups",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Cluster Backups.",
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
		ShortDesc:  "Get a Cluster Backup",
		Example:    getBackupExample,
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Backup by using its ID.\n\nRequired values to run command:\n\n* Backup Id",
		PreCmdRun:  PreRunBackupId,
		CmdRun:     RunBackupGet,
		InitClient: true,
	})
	get.AddStringFlag(constants.FlagBackupId, constants.FlagIdShort, "", "The unique ID of the Backup", core.RequiredFlagOption())
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

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	return c.Out(table.Sprint(allBackupCols, backups, cols, table.WithPrefix("items")))
}

func RunBackupGet(c *core.CommandConfig) error {
	c.Verbose("Backup ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId)))
	c.Verbose("Getting Backup...")

	backup, _, err := client.Must().PostgresClient.BackupsApi.ClustersBackupsFindById(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	return c.Out(table.Sprint(allBackupCols, backup, cols))
}

func ClusterBackupCmd() *core.Command {
	ctx := context.TODO()
	clusterBackupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "PostgreSQL Cluster Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster backup` allow you to list PostgreSQL Backups from a specific Cluster.",
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
		ShortDesc:  "List Cluster Backups from a Cluster",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Cluster Backups from a specific Cluster.",
		Example:    listBackupExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterBackupList,
		InitClient: true,
	})
	list.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", nil, table.ColsMessage(allBackupCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allBackupCols), cobra.ShellCompDirectiveNoFileComp
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
	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	return c.Out(table.Sprint(allBackupCols, backups, cols, table.WithPrefix("items")))
}
