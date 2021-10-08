package cloudapi_dbaas_pgsql

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-dbaas-pgsql/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
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
			Long:             "The sub-commands of `ionosctl dbaas-pgsql backup` allow you to create, list, get, update and delete PostgreSQL Backups.",
			TraverseChildren: true,
		},
	}
	globalFlags := backupCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultBackupCols, printer.ColsMessage(allBackupCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(backupCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = backupCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allBackupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, backupCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
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

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, backupCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
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
	get.AddStringFlag(cloudapidbaaspgsql.ArgBackupId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.BackupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return backupCmd
}

func PreRunBackupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapidbaaspgsql.ArgBackupId)
}

func RunBackupList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Backups...")
	backups, _, err := c.CloudApiDbaasPgsqlServices.Backups().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupPrint(c, getBackups(backups)))
}

func RunBackupGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Backup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgBackupId)))
	c.Printer.Verbose("Getting Backup...")
	backup, _, err := c.CloudApiDbaasPgsqlServices.Backups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgBackupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupPrint(c, []resources.ClusterBackup{*backup}))
}

func ClusterBackupCmd() *core.Command {
	ctx := context.TODO()
	clusterBackupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "PostgreSQL Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas-pgsql cluster backup` allow you to list PostgreSQL Backups from a specified Cluster.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterBackupCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql.cluster",
		Resource:   "backup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Cluster Backups from a Cluster",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Cluster Backups from a specified Cluster.",
		Example:    listBackupExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterBackupList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapidbaaspgsql.ArgClusterId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.ClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(config.ArgCols, "", defaultBackupCols, printer.ColsMessage(allBackupCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allBackupCols, cobra.ShellCompDirectiveNoFileComp
	})

	return clusterBackupCmd
}

func RunClusterBackupList(c *core.CommandConfig) error {
	c.Printer.Verbose("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId)))
	c.Printer.Verbose("Getting Backups from Cluster...")
	backups, _, err := c.CloudApiDbaasPgsqlServices.Backups().ListBackups(viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getBackupPrint(c, getBackups(backups)))
}

// Output Printing

var (
	defaultBackupCols = []string{"BackupId", "ClusterId", "DisplayName", "Type", "CreatedDate"}
	allBackupCols     = []string{"BackupId", "ClusterId", "DisplayName", "Type", "CreatedDate", "LastModifiedDate"}
)

type BackupPrint struct {
	BackupId         string `json:"BackupId,omitempty"`
	ClusterId        string `json:"ClusterId,omitempty"`
	DisplayName      string `json:"DisplayName,omitempty"`
	Type             string `json:"Type,omitempty"`
	CreatedDate      string `json:"CreatedDate,omitempty"`
	LastModifiedDate string `json:"LastModifiedDate,omitempty"`
}

func getBackupPrint(c *core.CommandConfig, dcs []resources.ClusterBackup) printer.Result {
	r := printer.Result{}
	if c != nil {
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getBackupsKVMaps(dcs)
			if strings.Contains(c.Namespace, "cluster") {
				r.Columns = getBackupCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getBackupCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
		}
	}
	return r
}

func getBackupCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultBackupCols
	}
	columnsMap := map[string]string{
		"BackupId":         "BackupId",
		"DisplayName":      "DisplayName",
		"ClusterId":        "ClusterId",
		"Type":             "Type",
		"CreatedDate":      "CreatedDate",
		"LastModifiedDate": "LastModifiedDate",
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

func getBackups(backups resources.ClusterBackupList) []resources.ClusterBackup {
	c := make([]resources.ClusterBackup, 0)
	if data, ok := backups.GetDataOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.ClusterBackup{ClusterBackup: d})
		}
	}
	return c
}

func getBackupsKVMaps(backups []resources.ClusterBackup) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(backups))
	for _, backup := range backups {
		var backupPrint BackupPrint
		if idOk, ok := backup.GetIdOk(); ok && idOk != nil {
			backupPrint.BackupId = *idOk
		}
		if displayNameOk, ok := backup.GetDisplayNameOk(); ok && displayNameOk != nil {
			backupPrint.DisplayName = *displayNameOk
		}
		if clusterIdOk, ok := backup.GetClusterIdOk(); ok && clusterIdOk != nil {
			backupPrint.ClusterId = *clusterIdOk
		}
		if typeOk, ok := backup.GetTypeOk(); ok && typeOk != nil {
			backupPrint.Type = *typeOk
		}
		if metadataOk, ok := backup.GetMetadataOk(); ok && metadataOk != nil {
			if createdDateOk, ok := metadataOk.GetCreatedDateOk(); ok && createdDateOk != nil {
				createdDate := *createdDateOk
				backupPrint.CreatedDate = createdDate.String()
			}
			if lastModifiedDateOk, ok := metadataOk.GetLastModifiedDateOk(); ok && lastModifiedDateOk != nil {
				lastModifiedDate := *lastModifiedDateOk
				backupPrint.LastModifiedDate = lastModifiedDate.String()
			}
		}
		o := structs.Map(backupPrint)
		out = append(out, o)
	}
	return out
}
