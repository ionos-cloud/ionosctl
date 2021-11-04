package pg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/structs"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/pg/completer"
	"github.com/ionos-cloud/ionosctl/commands/pg/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-pg"
	"github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl pg cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}
	globalFlags := clusterCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(clusterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List PostgreSQL Clusters",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:    listClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})
	list.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Cluster",
		Example:    getClusterExample,
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})
	get.AddStringFlag(dbaaspg.ArgClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dbaaspg.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-pgsql",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL Cluster",
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the VDC (VirtualDataCenter), the unique ID of the LAN. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the VDC.

Required values to run command:

* Datacenter Id
* Lan Id
* IP`,
		Example:    createClusterExample,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(dbaaspg.ArgVersion, dbaaspg.ArgVersionShort, "13", "The PostgreSQL version of your Cluster")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(dbaaspg.ArgReplicas, dbaaspg.ArgReplicasShort, 1, "The number of replicas in your cluster. Minimum: 1. Maximum: 5")
	create.AddIntFlag(dbaaspg.ArgCpuCoreCount, "", 4, "The number of CPU cores per replica")
	create.AddStringFlag(dbaaspg.ArgRamSize, "", "2Gi", "The amount of memory per replica in IEC format. Value must be a multiple of 1024Mi and at least 2048Mi")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgRamSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048Mi", "2Gi", "4Gi", "10Gi"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgSyncMode, dbaaspg.ArgSyncModeShort, "asynchronous", "Represents different modes of replication")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"asynchronous", "synchronous", "strictly_synchronous"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgStorageSize, "", "20Gi", "The amount of storage per replica. It is expected IEC format like 2Gi or 500Mi")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"20Gi", "500Mi", "2Gi", "50Gi"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgStorageType, "", "HDD", "The storage type used in your cluster")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLocation, "", "de/fra", "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests). If not set, it will be used VDC's location")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"de/fra", "de/txl", "gb/lhr"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddStringFlag(dbaaspg.ArgDatacenterId, dbaaspg.ArgDatacenterIdShort, "", "The unique ID of the VDC to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLanId, "", "", "The unique Lan ID", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, dbaaspg.ArgDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgIpAddress, "", "", "The private IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Example: 192.168.1.100/24", core.RequiredFlagOption())
	create.AddStringFlag(dbaaspg.ArgBackupId, dbaaspg.ArgBackupIdShort, "", "The unique ID of the backup you want to restore")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgTime, "", "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
	create.AddStringFlag(dbaaspg.ArgUsername, "", "db-admin", "Username for the database user to be created. Some system usernames are restricted (e.g. postgres, admin, standby)")
	create.AddStringFlag(dbaaspg.ArgPassword, "", "password", "Password for the database user to be created")
	create.AddStringFlag(dbaaspg.ArgMaintenanceTime, dbaaspg.ArgMaintenanceTimeShort, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59")
	create.AddStringFlag(dbaaspg.ArgMaintenanceDay, dbaaspg.ArgMaintenanceDayShort, "", "WeekDay for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-pgsql",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a PostgreSQL Cluster",
		LongDesc: `Use this command to update attributes of a PostgreSQL Cluster.

Required values to run command:

* Cluster Id`,
		Example:    updateClusterExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddStringFlag(dbaaspg.ArgClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgVersion, dbaaspg.ArgVersionShort, "", "The PostgreSQL version of your cluster")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(dbaaspg.ArgReplicas, dbaaspg.ArgReplicasShort, 0, "The number of replicas in your cluster")
	update.AddIntFlag(dbaaspg.ArgCpuCoreCount, "", 0, "The number of CPU cores per replica")
	update.AddStringFlag(dbaaspg.ArgRamSize, "", "", "The amount of memory per replica in IEC format. Value must be a multiple of 1024Mi and at least 2048Mi")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgRamSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2Gi", "4Gi", "2048Mi", "10Gi"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgStorageSize, "", "", "The amount of storage per replica. It is expected IEC format like 2Gi or 500Mi")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"20Gi", "500Mi", "2Gi", "50Gi"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(dbaaspg.ArgMaintenanceTime, dbaaspg.ArgMaintenanceTimeShort, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59")
	update.AddStringFlag(dbaaspg.ArgMaintenanceDay, dbaaspg.ArgMaintenanceDayShort, "", "WeekDay for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Restore Command
	*/
	restoreCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-pgsql",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a PostgreSQL Cluster",
		LongDesc: `Use this command to trigger an in-place restore of the specified PostgreSQL Cluster.

Required values to run command:

* Cluster Id
* Backup Id`,
		Example:    restoreClusterExample,
		PreCmdRun:  PreRunClusterBackupIds,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddStringFlag(dbaaspg.ArgClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(dbaaspg.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(dbaaspg.ArgBackupId, "", "", "The unique ID of the backup you want to restore", core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(dbaaspg.ArgBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(dbaaspg.ArgTime, "", "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
	restoreCmd.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	restoreCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-pgsql",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a PostgreSQL Cluster",
		LongDesc: `Use this command to delete a specified PostgreSQL Cluster from your account.

Required values to run command:

* Cluster Id`,
		Example:    deleteClusterExample,
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(dbaaspg.ArgClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(dbaaspg.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgAll, config.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Delete all Clusters after filtering based on name. Can be used with --all flag")

	clusterCmd.AddCommand(ClusterBackupCmd())

	return clusterCmd
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgClusterId)
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{dbaaspg.ArgClusterId}, []string{config.ArgAll})
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		if !viper.IsSet(core.GetFlagName(c.NS, config.ArgAll)) {
			return errors.New("error: name flag can to be used with the --all flag")
		}
	}
	return nil
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgDatacenterId, dbaaspg.ArgLanId, dbaaspg.ArgIpAddress)
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgTime)) {
		if !viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)) {
			return errors.New("error: recovery target time can be used with --backup-id flag")
		}
	}
	return nil
}

func PreRunClusterBackupIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgClusterId, dbaaspg.ArgBackupId)
}

func RunClusterList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Clusters...")
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		c.Printer.Verbose("Filtering after Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	}
	clusters, _, err := c.CloudApiDbaasPgsqlServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, getClusters(clusters)))
}

func RunClusterGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId)))
	c.Printer.Verbose("Getting Cluster...")
	if err := utils.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))); err != nil {
		return err
	}
	cluster, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, []resources.Cluster{*cluster}))
}

func RunClusterCreate(c *core.CommandConfig) error {
	var recoveryTargetTime time.Time
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)) {
		c.Printer.Verbose("Backup ID: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgTime)) {
		c.Printer.Verbose("RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgTime)))
		recoveryTargetTime, err = time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgTime)))
		if err != nil {
			return err
		}
	}
	c.Printer.Verbose("Creating Cluster...")
	cluster, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Create(*input,
		viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)), recoveryTargetTime)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := cluster.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.ClusterStateInterrogator, *id); err != nil {
				return err
			}
			if cluster, _, err = c.CloudApiDbaasPgsqlServices.Clusters().Get(*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new Cluster Id")
		}
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.Cluster{*cluster}))
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))
	c.Printer.Verbose("Cluster ID: %v", clusterId)
	input := getPatchClusterRequest(c)
	c.Printer.Verbose("Updating Cluster...")
	item, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Update(clusterId, input)
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.Cluster{*item}))
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))
	backupId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId))
	c.Printer.Verbose("Cluster ID: %v", clusterId)
	c.Printer.Verbose("Backup ID: %v", backupId)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "restore cluster"); err != nil {
		return err
	}
	input := resources.CreateRestoreRequest{
		CreateRestoreRequest: sdkgo.CreateRestoreRequest{
			BackupId: &backupId,
		},
	}
	if viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgTime)) != "" {
		c.Printer.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgTime)))
		recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgTime)))
		if err != nil {
			return err
		}
		input.SetRecoveryTargetTime(recoveryTargetTime)
	}
	c.Printer.Verbose("Restoring Cluster from Backup...")
	resp, err := c.CloudApiDbaasPgsqlServices.Restores().Restore(clusterId, input)
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))); err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(resp, c, nil))
}

func RunClusterDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgAll)) {
		resp, err = ClusterDeleteAll(c)
		if err != nil {
			return err
		}
	} else {
		clusterId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))
		c.Printer.Verbose("Cluster ID: %v", clusterId)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete cluster"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting Cluster...")
		resp, err = c.CloudApiDbaasPgsqlServices.Clusters().Delete(clusterId)
		if err != nil {
			return err
		}
	}
	return c.Printer.Print(getClusterPrint(resp, c, nil))
}

func ClusterDeleteAll(c *core.CommandConfig) (resp *resources.Response, err error) {
	c.Printer.Verbose("Getting all Clusters...")
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		c.Printer.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	}
	clusters, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	if err != nil {
		return nil, err
	}
	if dataOk, ok := clusters.GetDataOk(); ok && dataOk != nil {
		for _, cluster := range *dataOk {
			var log string
			if nameOk, ok := cluster.GetDisplayNameOk(); ok && nameOk != nil {
				log = fmt.Sprintf("Cluster Name: %s", *nameOk)
			}
			if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
				log = fmt.Sprintf("%s; Cluster Id: %s", log, *idOk)
			}
			c.Printer.Verbose(log)
		}
	}
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all clusters"); err != nil {
		return nil, err
	}
	if dataOk, ok := clusters.GetDataOk(); ok && dataOk != nil {
		for _, cluster := range *dataOk {
			if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
				c.Printer.Verbose("Cluster ID: %v", *idOk)
				c.Printer.Verbose("Deleting Cluster...")
				resp, err = c.CloudApiDbaasPgsqlServices.Clusters().Delete(*idOk)
				if err != nil {
					return resp, err
				}
			}
		}
	}
	return resp, err
}

func getCreateClusterRequest(c *core.CommandConfig) (*resources.CreateClusterRequest, error) {
	input := resources.CreateClusterRequest{}
	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgVersion))
	c.Printer.Verbose("PostgresVersion: %v", pgsqlVersion)
	input.SetPostgresVersion(pgsqlVersion)
	syncMode := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgSyncMode))
	c.Printer.Verbose("SynchronizationMode: %v", syncMode)
	input.SetSynchronizationMode(sdkgo.SynchronizationMode(syncMode))
	replicas := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgReplicas))
	c.Printer.Verbose("Replicas: %v", replicas)
	input.SetReplicas(replicas)
	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgCpuCoreCount))
	c.Printer.Verbose("CpuCoreCount: %v", cpuCoreCount)
	input.SetCpuCoreCount(cpuCoreCount)
	ramSize := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRamSize))
	c.Printer.Verbose("RamSize: %v", ramSize)
	input.SetRamSize(ramSize)
	storageSize := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize))
	c.Printer.Verbose("StorageSize: %v", storageSize)
	input.SetStorageSize(storageSize)
	storageType := sdkgo.StorageType(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageType)))
	c.Printer.Verbose("StorageType: %v", storageType)
	input.SetStorageType(storageType)
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLocation)) {
		location := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLocation))
		c.Printer.Verbose("Location: %v", location)
		input.SetLocation(location)
	} else {
		c.Printer.Verbose("Getting Location from VDC...")
		vdc, _, err := c.CloudApiV6Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId)))
		if err != nil {
			return nil, err
		}
		if properties, ok := vdc.GetPropertiesOk(); ok && properties != nil {
			if location, ok := properties.GetLocationOk(); ok && location != nil {
				c.Printer.Verbose("Location: %v", *location)
				input.SetLocation(*location)
			}
		}
	}
	displayName := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))
	c.Printer.Verbose("DisplayName: %v", displayName)
	input.SetDisplayName(displayName)
	dbuser := sdkgo.DBUser{}
	username := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgUsername))
	c.Printer.Verbose("DBUser - Username: %v", username)
	dbuser.SetUsername(username)
	password := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgPassword))
	c.Printer.Verbose("DBUser - Password: %v", password)
	dbuser.SetPassword(password)
	input.SetCredentials(dbuser)
	vdcConnection := sdkgo.VDCConnection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId))
	c.Printer.Verbose("VDCConnection - VdcId: %v", vdcId)
	vdcConnection.SetVdcId(vdcId)
	lanId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLanId))
	c.Printer.Verbose("VDCConnection - LanId: %v", lanId)
	vdcConnection.SetLanId(lanId)
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgIpAddress)) {
		ip := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgIpAddress))
		c.Printer.Verbose("VDCConnection - IpAddress: %v", ip)
		vdcConnection.SetIpAddress(ip)
	}
	input.SetVdcConnections([]sdkgo.VDCConnection{vdcConnection})
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay)) {
		maintenanceWindow := sdkgo.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime))
			c.Printer.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay))
			c.Printer.Verbose("MaintenanceWindow - WeekDay: %v", maintenanceDay)
			maintenanceWindow.SetWeekday(maintenanceDay)
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}
	return &input, nil
}

func getPatchClusterRequest(c *core.CommandConfig) resources.PatchClusterRequest {
	input := resources.PatchClusterRequest{}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCpuCoreCount)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgCpuCoreCount))
		c.Printer.Verbose("CpuCoreCount: %v", cpuCoreCount)
		input.SetCpuCoreCount(cpuCoreCount)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgReplicas)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgReplicas))
		c.Printer.Verbose("Replicas: %v", replicas)
		input.SetReplicas(replicas)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRamSize)) {
		ramSize := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRamSize))
		c.Printer.Verbose("RamSize: %v", ramSize)
		input.SetRamSize(ramSize)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize)) {
		storageSize := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize))
		c.Printer.Verbose("StorageSize: %v", storageSize)
		input.SetStorageSize(storageSize)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgVersion)) {
		pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgVersion))
		c.Printer.Verbose("PostgresVersion: %v", pgsqlVersion)
		input.SetPostgresVersion(pgsqlVersion)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))
		c.Printer.Verbose("DisplayName: %v", displayName)
		input.SetDisplayName(displayName)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime)) || viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay)) {
		maintenanceWindow := sdkgo.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceTime))
			c.Printer.Verbose("MaintenanceTime: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgMaintenanceDay))
			c.Printer.Verbose("MaintenanceWeekDay: %v", maintenanceDay)
			maintenanceWindow.SetWeekday(maintenanceDay)
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}
	return input
}

// Output Printing

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "Location", "DatacenterId", "LanId", "IpAddress", "Replicas", "LifecycleStatus"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "Location", "BackupEnabled", "LifecycleStatus", "PostgresVersion", "Replicas", "RamSize", "CpuCoreCount",
		"StorageSize", "StorageType", "DatacenterId", "LanId", "IpAddress", "MaintenanceWindow"}
)

type ClusterPrint struct {
	ClusterId         string `json:"ClusterId,omitempty"`
	Location          string `json:"Location,omitempty"`
	BackupEnabled     bool   `json:"BackupEnabled,omitempty"`
	LifecycleStatus   string `json:"LifecycleStatus,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty"`
	PostgresVersion   string `json:"PostgresVersion,omitempty"`
	Replicas          int32  `json:"Replicas,omitempty"`
	RamSize           string `json:"RamSize,omitempty"`
	CpuCoreCount      int32  `json:"CpuCoreCount,omitempty"`
	StorageSize       string `json:"StorageSize,omitempty"`
	StorageType       string `json:"StorageType,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	LanId             string `json:"LanId,omitempty"`
	IpAddress         string `json:"IpAddress,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
}

func getClusterPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.Cluster) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getClustersKVMaps(dcs)
			r.Columns = getClusterCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getClusterCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultClusterCols
	}

	columnsMap := map[string]string{
		"ClusterId":         "ClusterId",
		"DisplayName":       "DisplayName",
		"Location":          "Location",
		"PostgresVersion":   "PostgresVersion",
		"BackupEnabled":     "BackupEnabled",
		"LifecycleStatus":   "LifecycleStatus",
		"Replicas":          "Replicas",
		"CpuCoreCount":      "CpuCoreCount",
		"StorageSize":       "StorageSize",
		"StorageType":       "StorageType",
		"DatacenterId":      "DatacenterId",
		"LanId":             "LanId",
		"IpAddress":         "IpAddress",
		"MaintenanceWindow": "MaintenanceWindow",
	}
	var clusterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			clusterCols = append(clusterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return clusterCols
}

func getClusters(clusters resources.ClusterList) []resources.Cluster {
	c := make([]resources.Cluster, 0)
	if data, ok := clusters.GetDataOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.Cluster{Cluster: d})
		}
	}
	return c
}

func getClustersKVMaps(clusters []resources.Cluster) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if displayNameOk, ok := cluster.GetDisplayNameOk(); ok && displayNameOk != nil {
			clusterPrint.DisplayName = *displayNameOk
		}
		if locationOk, ok := cluster.GetLocationOk(); ok && locationOk != nil {
			clusterPrint.Location = *locationOk
		}
		if backupEnabledOk, ok := cluster.GetBackupEnabledOk(); ok && backupEnabledOk != nil {
			clusterPrint.BackupEnabled = *backupEnabledOk
		}
		if vdcConnectionsOk, ok := cluster.GetVdcConnectionsOk(); ok && vdcConnectionsOk != nil {
			for _, vdcConnection := range *vdcConnectionsOk {
				if vdcIdOk, ok := vdcConnection.GetVdcIdOk(); ok && vdcIdOk != nil {
					clusterPrint.DatacenterId = *vdcIdOk
				}
				if lanIdOk, ok := vdcConnection.GetLanIdOk(); ok && lanIdOk != nil {
					clusterPrint.LanId = *lanIdOk
				}
				if ipAddressOk, ok := vdcConnection.GetIpAddressOk(); ok && ipAddressOk != nil {
					clusterPrint.IpAddress = *ipAddressOk
				}
			}
		}
		if lifecycleStatusOk, ok := cluster.GetLifecycleStatusOk(); ok && lifecycleStatusOk != nil {
			clusterPrint.LifecycleStatus = *lifecycleStatusOk
		}
		if postgresVersionOk, ok := cluster.GetPostgresVersionOk(); ok && postgresVersionOk != nil {
			clusterPrint.PostgresVersion = *postgresVersionOk
		}
		if replicasOk, ok := cluster.GetReplicasOk(); ok && replicasOk != nil {
			clusterPrint.Replicas = *replicasOk
		}
		if ramSizeOk, ok := cluster.GetRamSizeOk(); ok && ramSizeOk != nil {
			clusterPrint.RamSize = *ramSizeOk
		}
		if cpuCoreCountOk, ok := cluster.GetCpuCoreCountOk(); ok && cpuCoreCountOk != nil {
			clusterPrint.CpuCoreCount = *cpuCoreCountOk
		}
		if storageSizeOk, ok := cluster.GetStorageSizeOk(); ok && storageSizeOk != nil {
			clusterPrint.StorageSize = *storageSizeOk
		}
		if storageTypeOk, ok := cluster.GetStorageTypeOk(); ok && storageTypeOk != nil {
			clusterPrint.StorageType = string(*storageTypeOk)
		}
		if maintenanceWindowOk, ok := cluster.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
			var maintenanceWindow string
			if weekdayOk, ok := maintenanceWindowOk.GetWeekdayOk(); ok && weekdayOk != nil {
				maintenanceWindow = *weekdayOk
			}
			if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
				maintenanceWindow = fmt.Sprintf("%s %s", maintenanceWindow, *timeOk)
			}
			clusterPrint.MaintenanceWindow = maintenanceWindow
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
