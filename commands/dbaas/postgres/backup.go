package postgres

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultBackupCols, printer.ColsMessage(allBackupCols))
	_ = viper.BindPFlag(core.GetFlagName(backupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = backupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allBackupCols, cobra.ShellCompDirectiveNoFileComp
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
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	get.AddStringFlag(dbaaspg.ArgBackupId, dbaaspg.ArgIdShort, "", dbaaspg.BackupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dbaaspg.ArgBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return backupCmd
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgBackupId)
}

func RunBackupList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Backups..."))

	backups, _, err := c.CloudApiDbaasPgsqlServices.Backups().List()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("items", allBackupJSONPaths, backups.ClusterBackupList,
		printer.GetHeaders(allBackupCols, defaultBackupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)
	return nil
}

func RunBackupGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Backup ID: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId))))
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Backup..."))

	backup, _, err := c.CloudApiDbaasPgsqlServices.Backups().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)))
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", allBackupJSONPaths, backup.BackupResponse,
		printer.GetHeaders(allBackupCols, defaultBackupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)
	return nil
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
	list.AddUUIDFlag(constants.FlagClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", defaultBackupCols, printer.ColsMessage(allBackupCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allBackupCols, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return clusterBackupCmd
}

func RunClusterBackupList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Backups from Cluster..."))

	backups, _, err := c.CloudApiDbaasPgsqlServices.Backups().ListBackups(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if err != nil {
		return err
	}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("items", allBackupJSONPaths, backups.ClusterBackupList,
		printer.GetHeaders(allBackupCols, defaultBackupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)
	return nil
}

// Output Printing

var (
	allBackupJSONPaths = map[string]string{
		"BackupId":                   "id",
		"ClusterId":                  "properties.clusterId",
		"EarliestRecoveryTargetTime": "properties.earliestRecoveryTargetTime",
		"Version":                    "properties.version",
		"Active":                     "properties.active",
		"CreatedDate":                "metadata.createdDate",
		"State":                      "metadata.state",
	}

	defaultBackupCols = []string{"BackupId", "ClusterId", "CreatedDate", "EarliestRecoveryTargetTime", "Active", "State"}
	allBackupCols     = []string{"BackupId", "ClusterId", "Active", "CreatedDate", "EarliestRecoveryTargetTime", "Version", "State"}
)

func getBackupCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultBackupCols
	}
	columnsMap := map[string]string{
		"BackupId":                   "BackupId",
		"ClusterId":                  "ClusterId",
		"EarliestRecoveryTargetTime": "EarliestRecoveryTargetTime",
		"CreatedDate":                "CreatedDate",
		"Version":                    "Version",
		"Active":                     "Active",
		"State":                      "State",
	}
	var backupCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			backupCols = append(backupCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return backupCols
}
