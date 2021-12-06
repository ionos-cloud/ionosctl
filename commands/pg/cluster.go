package pg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
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
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the Datacenter.

Required values to run command:

* Datacenter Id
* Lan Id
* IP
* Credentials for the User: Username and Password`,
		Example:    createClusterExample,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(dbaaspg.ArgVersion, dbaaspg.ArgVersionShort, "13", "The PostgreSQL version of your Cluster")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(dbaaspg.ArgInstances, dbaaspg.ArgInstancesShort, 1, "The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5")
	create.AddIntFlag(dbaaspg.ArgCores, "", 2, "The number of CPU cores per instance. Minimum: 1")
	create.AddStringFlag(dbaaspg.ArgRam, "", "3GB", "The amount of memory per instance. Size must be specified in multiples of 256. Minimum: 2048. The default unit is MB. e.g. --ram 2048 or --ram 2048MB")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgSyncMode, dbaaspg.ArgSyncModeShort, "ASYNCHRONOUS", "Represents different modes of replication")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgStorageSize, "", "20GB", "The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgStorageType, "", "HDD", "The storage type used in your cluster")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLocation, "", "", "The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation. If not set, it will be used Datacenter's location")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddStringFlag(dbaaspg.ArgDatacenterId, dbaaspg.ArgDatacenterIdShort, "", "The unique ID of the Datacenter to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLanId, "", "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, dbaaspg.ArgDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgCidr, "", "", "The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Example: 192.168.1.100/24", core.RequiredFlagOption())
	create.AddStringFlag(dbaaspg.ArgBackupId, dbaaspg.ArgBackupIdShort, "", "The unique ID of the backup you want to restore")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgRecoveryTime, dbaaspg.ArgRecoveryTimeShort, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
	create.AddStringFlag(dbaaspg.ArgUsername, dbaaspg.ArgUsernameShort, "db-admin", "Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby)", core.RequiredFlagOption())
	create.AddStringFlag(dbaaspg.ArgPassword, dbaaspg.ArgPasswordShort, "", "Password for the initial postgres user", core.RequiredFlagOption())
	create.AddStringFlag(dbaaspg.ArgMaintenanceTime, dbaaspg.ArgMaintenanceTimeShort, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59")
	create.AddStringFlag(dbaaspg.ArgMaintenanceDay, dbaaspg.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
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
	update.AddIntFlag(dbaaspg.ArgInstances, dbaaspg.ArgInstancesShort, 0, "The number of instances in your cluster")
	update.AddIntFlag(dbaaspg.ArgCores, "", 0, "The number of CPU cores per instance")
	update.AddStringFlag(dbaaspg.ArgRam, "", "", "The amount of memory per instance. Size must be specified in multiples of 256. Minimum: 2048. The default unit is MB. e.g. --ram 2048 or --ram 2048MB")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgStorageSize, "", "", "The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(dbaaspg.ArgMaintenanceTime, dbaaspg.ArgMaintenanceTimeShort, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Example: 16:30:59")
	update.AddStringFlag(dbaaspg.ArgMaintenanceDay, dbaaspg.ArgMaintenanceDayShort, "", "Day of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")

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
		return completer.BackupsIdsForCluster(os.Stderr, viper.GetString(core.GetFlagName(restoreCmd.NS, dbaaspg.ArgClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(dbaaspg.ArgRecoveryTime, dbaaspg.ArgRecoveryTimeShort, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
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
	deleteCmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")

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
	err := core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgDatacenterId, dbaaspg.ArgLanId, dbaaspg.ArgCidr, dbaaspg.ArgPassword)
	if err != nil {
		return err
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
	return c.Printer.Print(getClusterPrint(nil, c, []resources.ClusterResponse{*cluster}))
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Creating Cluster...")
	cluster, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Create(*input)
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
	return c.Printer.Print(getClusterPrint(resp, c, []resources.ClusterResponse{*cluster}))
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))
	c.Printer.Verbose("Cluster ID: %v", clusterId)
	input, err := getPatchClusterRequest(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating Cluster...")
	item, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Update(clusterId, *input)
	if err != nil {
		return err
	}
	if err = utils.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgClusterId))); err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.ClusterResponse{*item}))
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
	if viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) != "" {
		c.Printer.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)))
		recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)))
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
	if dataOk, ok := clusters.GetItemsOk(); ok && dataOk != nil {
		for _, cluster := range *dataOk {
			var log string
			if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
				if nameOk, ok := propertiesOk.GetDisplayNameOk(); ok && nameOk != nil {
					log = fmt.Sprintf("Cluster Name: %s", *nameOk)
				}
			}
			if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
				log = fmt.Sprintf("%s; Cluster ID: %s", log, *idOk)
			}
			c.Printer.Print(log)
		}
	}
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all clusters"); err != nil {
		return nil, err
	}
	if dataOk, ok := clusters.GetItemsOk(); ok && dataOk != nil {
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
	inputCluster := resources.CreateClusterRequest{}
	input := sdkgo.CreateClusterProperties{}
	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgVersion))
	c.Printer.Verbose("PostgresVersion: %v", pgsqlVersion)
	input.SetPostgresVersion(pgsqlVersion)
	syncMode := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgSyncMode)))
	c.Printer.Verbose("SynchronizationMode: %v", syncMode)
	input.SetSynchronizationMode(sdkgo.SynchronizationMode(syncMode))
	replicas := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgInstances))
	c.Printer.Verbose("Instances: %v", replicas)
	input.SetInstances(replicas)
	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgCores))
	c.Printer.Verbose("Cores: %v", cpuCoreCount)
	input.SetCores(cpuCoreCount)
	// Convert Ram
	size, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRam)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}
	input.SetRam(int32(size))
	c.Printer.Verbose("Ram: %v[MB]", int32(size))
	// Convert StorageSize
	storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize)), utils.MegaBytes)
	if err != nil {
		return nil, err
	}
	input.SetStorageSize(int32(storageSize))
	c.Printer.Verbose("StorageSize: %v[MB]", int32(storageSize))
	storageType := sdkgo.StorageType(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageType)))
	c.Printer.Verbose("StorageType: %v", storageType)
	input.SetStorageType(storageType)
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLocation)) {
		location := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLocation))
		c.Printer.Verbose("Location: %v", location)
		input.SetLocation(sdkgo.Location(location))
	} else {
		c.Printer.Verbose("Getting Location from VDC...")
		vdc, _, err := c.CloudApiV6Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId)))
		if err != nil {
			return nil, err
		}
		if properties, ok := vdc.GetPropertiesOk(); ok && properties != nil {
			if location, ok := properties.GetLocationOk(); ok && location != nil {
				c.Printer.Verbose("Location: %v", *location)
				input.SetLocation(sdkgo.Location(*location))
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
	vdcConnection := sdkgo.Connection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId))
	c.Printer.Verbose("Connection - DatacenterId: %v", vdcId)
	vdcConnection.SetDatacenterId(vdcId)
	lanId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLanId))
	c.Printer.Verbose("Connection - LanId: %v", lanId)
	vdcConnection.SetLanId(lanId)
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCidr)) {
		ip := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgCidr))
		c.Printer.Verbose("Connection - Cidr: %v", ip)
		vdcConnection.SetCidr(ip)
	}
	input.SetConnections([]sdkgo.Connection{vdcConnection})
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
			c.Printer.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(sdkgo.DayOfTheWeek(maintenanceDay))
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)) ||
		viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) {
		createRestoreRequest := sdkgo.CreateRestoreRequest{}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) {
			recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)))
			if err != nil {
				return nil, err
			}
			c.Printer.Verbose("From Backup - RecoveryTargetTime [RFC3339 format]: %v", recoveryTargetTime)
			createRestoreRequest.SetRecoveryTargetTime(recoveryTargetTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgBackupId)) {
			backupId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgBackupId))
			c.Printer.Verbose("From Backup - BackupId: %v", backupId)
			createRestoreRequest.SetBackupId(backupId)
		}
		input.SetFromBackup(createRestoreRequest)
	}
	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

func getPatchClusterRequest(c *core.CommandConfig) (*resources.PatchClusterRequest, error) {
	inputCluster := resources.PatchClusterRequest{}
	input := sdkgo.PatchClusterProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCores)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgCores))
		c.Printer.Verbose("Cores: %v", cpuCoreCount)
		input.SetCores(cpuCoreCount)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgInstances)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, dbaaspg.ArgInstances))
		c.Printer.Verbose("Instances: %v", replicas)
		input.SetInstances(replicas)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRam)) {
		// Convert Ram
		size, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRam)), utils.MegaBytes)
		if err != nil {
			return nil, err
		}
		input.SetRam(int32(size))
		c.Printer.Verbose("Ram: %vMB", int32(size))
	}
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize)) {
		// Convert StorageSize
		storageSize, err := utils.ConvertSize(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgStorageSize)), utils.MegaBytes)
		if err != nil {
			return nil, err
		}
		input.SetStorageSize(int32(storageSize))
		c.Printer.Verbose("StorageSize: %vMB", storageSize)
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
			c.Printer.Verbose("MaintenanceDayOfWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(sdkgo.DayOfTheWeek(maintenanceDay))
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}
	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

// Output Printing

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "Location", "DatacenterId", "LanId", "Cidr", "Instances", "State"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "Location", "State", "PostgresVersion", "Instances", "Ram", "Cores",
		"StorageSize", "StorageType", "DatacenterId", "LanId", "Cidr", "MaintenanceWindow"}
)

type ClusterPrint struct {
	ClusterId         string `json:"ClusterId,omitempty"`
	Location          string `json:"Location,omitempty"`
	State             string `json:"State,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty"`
	PostgresVersion   string `json:"PostgresVersion,omitempty"`
	Instances         int32  `json:"Instances,omitempty"`
	Ram               string `json:"Ram,omitempty"`
	Cores             int32  `json:"Cores,omitempty"`
	StorageSize       string `json:"StorageSize,omitempty"`
	StorageType       string `json:"StorageType,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	LanId             string `json:"LanId,omitempty"`
	Cidr              string `json:"Cidr,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
}

func getClusterPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.ClusterResponse) printer.Result {
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
		"State":             "State",
		"Ram":               "Ram",
		"Instances":         "Instances",
		"Cores":             "Cores",
		"StorageSize":       "StorageSize",
		"StorageType":       "StorageType",
		"DatacenterId":      "DatacenterId",
		"LanId":             "LanId",
		"Cidr":              "Cidr",
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

func getClusters(clusters resources.ClusterList) []resources.ClusterResponse {
	c := make([]resources.ClusterResponse, 0)
	if data, ok := clusters.GetItemsOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.ClusterResponse{ClusterResponse: d})
		}
	}
	return c
}

func getClustersKVMaps(clusters []resources.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetDisplayNameOk(); ok && displayNameOk != nil {
				clusterPrint.DisplayName = *displayNameOk
			}
			if locationOk, ok := propertiesOk.GetLocationOk(); ok && locationOk != nil {
				clusterPrint.Location = string(*locationOk)
			}
			if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
				for _, vdcConnection := range *vdcConnectionsOk {
					if vdcIdOk, ok := vdcConnection.GetDatacenterIdOk(); ok && vdcIdOk != nil {
						clusterPrint.DatacenterId = *vdcIdOk
					}
					if lanIdOk, ok := vdcConnection.GetLanIdOk(); ok && lanIdOk != nil {
						clusterPrint.LanId = *lanIdOk
					}
					if ipAddressOk, ok := vdcConnection.GetCidrOk(); ok && ipAddressOk != nil {
						clusterPrint.Cidr = *ipAddressOk
					}
				}
			}
			if postgresVersionOk, ok := propertiesOk.GetPostgresVersionOk(); ok && postgresVersionOk != nil {
				clusterPrint.PostgresVersion = *postgresVersionOk
			}
			if replicasOk, ok := propertiesOk.GetInstancesOk(); ok && replicasOk != nil {
				clusterPrint.Instances = *replicasOk
			}
			if ramSizeOk, ok := propertiesOk.GetRamOk(); ok && ramSizeOk != nil {
				clusterPrint.Ram = fmt.Sprintf("%vMB", *ramSizeOk)
			}
			if cpuCoreCountOk, ok := propertiesOk.GetCoresOk(); ok && cpuCoreCountOk != nil {
				clusterPrint.Cores = *cpuCoreCountOk
			}
			if storageSizeOk, ok := propertiesOk.GetStorageSizeOk(); ok && storageSizeOk != nil {
				clusterPrint.StorageSize = fmt.Sprintf("%vGB", *storageSizeOk)
			}
			if storageTypeOk, ok := propertiesOk.GetStorageTypeOk(); ok && storageTypeOk != nil {
				clusterPrint.StorageType = string(*storageTypeOk)
			}
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				var maintenanceWindow string
				if weekdayOk, ok := maintenanceWindowOk.GetDayOfTheWeekOk(); ok && weekdayOk != nil {
					maintenanceWindow = string(*weekdayOk)
				}
				if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
					maintenanceWindow = fmt.Sprintf("%s %s", maintenanceWindow, *timeOk)
				}
				clusterPrint.MaintenanceWindow = maintenanceWindow
			}
		}
		if metadataOk, ok := cluster.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				clusterPrint.State = string(*stateOk)
			}
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
