package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
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
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
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
	list.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
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
	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state [seconds]")
	get.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = get.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL Cluster",
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the Datacenter.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* Credentials for the database user: Username and Password`,
		Example:    createClusterExample,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "13", "The PostgreSQL version of your Cluster")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 1, "The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5")
	create.AddIntFlag(constants.FlagCores, "", 2, "The number of CPU cores per instance. Minimum: 1")
	create.AddStringFlag(constants.FlagRam, "", "3GB", "The amount of memory per instance. Size must be specified in multiples of 1024. The default unit is MB. Minimum: 2048. e.g. --ram 2048, --ram 2048MB, --ram 2GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "", "The S3 location where the backups will be stored")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"de", "eu-south-2", "eu-central-2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagSyncMode, constants.FlagSyncModeShort, "ASYNCHRONOUS", "Synchronization Mode. Represents different modes of replication")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "20GB", "The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageType, "", "HDD", "The storage type used in your cluster: HDD, SSD, SSD_PREMIUM, SSD_STANDARD. (Value \"SSD\" is deprecated. Use the equivalent \"SSD_PREMIUM\" instead)")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLocation, "", "", "The physical location where the cluster will be created. It cannot be modified after datacenter creation. If not set, it will be used Datacenter's location")
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddUUIDFlag(dbaaspg.ArgDatacenterId, dbaaspg.ArgDatacenterIdShort, "", "The unique ID of the Datacenter to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgLanId, dbaaspg.ArgLanIdShort, "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, dbaaspg.ArgDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgCidr, dbaaspg.ArgCidrShort, "", "The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. e.g.: 192.168.1.100/24", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagBackupId, "b", "", "The unique ID of the backup you want to restore")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dbaaspg.ArgRecoveryTime, dbaaspg.ArgRecoveryTimeShort, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
	create.AddStringFlag(dbaaspg.ArgDbUsername, dbaaspg.ArgDbUsernameShort, "", "Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby)", core.RequiredFlagOption())
	create.AddStringFlag(dbaaspg.ArgDbPassword, dbaaspg.ArgDbPasswordShort, "", "Password for the initial postgres user", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59")
	create.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	create.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = create.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
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
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "", "The PostgreSQL version of your cluster")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(dbaaspg.ArgRemoveConnection, "", false, "Remove the connection completely")
	update.AddUUIDFlag(dbaaspg.ArgDatacenterId, dbaaspg.ArgDatacenterIdShort, "", "The unique ID of the Datacenter to connect to your cluster. It has to be in the same location as the current datacenter")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgLanId, dbaaspg.ArgLanIdShort, "", "The unique ID of the LAN to connect your cluster to")
	_ = update.Command.RegisterFlagCompletionFunc(dbaaspg.ArgLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(update.NS, dbaaspg.ArgDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgCidr, dbaaspg.ArgCidrShort, "", "The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. e.g.: 192.168.1.100/24")
	update.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 0, "The number of instances in your cluster. Minimum: 0. Maximum: 5")
	update.AddIntFlag(constants.FlagCores, "", 0, "The number of CPU cores per instance")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory per instance. Size must be specified in multiples of 1024. The default unit is MB. Minimum: 2048. e.g. --ram 2048, --ram 2048MB, --ram 2GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2GB", "3GB", "4GB", "5GB", "10GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59")
	update.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "", "Day of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	update.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = update.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Restore Command
	*/
	restoreCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
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
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(constants.FlagBackupId, "", "", "The unique ID of the backup you want to restore", core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIdsForCluster(viper.GetString(core.GetFlagName(restoreCmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(dbaaspg.ArgRecoveryTime, dbaaspg.ArgRecoveryTimeShort, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")
	restoreCmd.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	restoreCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	restoreCmd.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a PostgreSQL Cluster",
		LongDesc: `Use this command to delete a specified PostgreSQL Cluster from your account. You can wait for the cluster to be deleted with the wait-for-deletion option.

Required values to run command:

* Cluster Id`,
		Example:    deleteClusterExample,
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")
	deleteCmd.AddBoolFlag(constants.ArgWaitForDelete, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be completely removed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, dbaaspg.DefaultClusterTimeout, "Timeout option for Cluster to be completely removed[seconds]")
	deleteCmd.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	clusterCmd.AddCommand(ClusterBackupCmd())

	return clusterCmd
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagClusterId}, []string{constants.ArgAll})
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
			return errors.New("error: name flag can to be used with the --all flag")
		}
	}
	return nil
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, dbaaspg.ArgDatacenterId, dbaaspg.ArgLanId, dbaaspg.ArgCidr, dbaaspg.ArgDbUsername, dbaaspg.ArgDbPassword)
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores)) < 1 {
		return errors.New("cores must be set to minimum: 1")
	}
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)) < 1 || viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)) > 5 {
		return errors.New("instances must be set to minimum: 1, maximum: 5")
	}
	return nil
}

func PreRunClusterBackupIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagBackupId)
}

func RunClusterList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Clusters..."))

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Filtering after Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))))
	}

	clusters, _, err := c.CloudApiDbaasPgsqlServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	clustersConverted, err := resource2table.ConvertDbaasPostgresClustersToTable(clusters.ClusterList)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(clusters, clustersConverted,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunClusterGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Cluster..."))

	if err := waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
		return err
	}

	cluster, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	clusterConverted, err := resource2table.ConvertDbaasPostgresClusterToTable(cluster.ClusterResponse)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(cluster, clusterConverted,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating Cluster..."))

	cluster, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Create(*input)
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := cluster.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, *id); err != nil {
				return err
			}

			if cluster, _, err = c.CloudApiDbaasPgsqlServices.Clusters().Get(*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new Cluster Id")
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	clusterConverted, err := resource2table.ConvertDbaasPostgresClusterToTable(cluster.ClusterResponse)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(cluster, clusterConverted,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))

	input, err := getPatchClusterRequest(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating Cluster..."))

	item, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Update(clusterId, *input)
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Wait 10 seconds before checking state..."))
		// TODO: Sleeping 10 seconds to make sure the cluster is in BUSY state. This will be removed in future releases.
		time.Sleep(10 * time.Second)

		if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	clusterConverted, err := resource2table.ConvertDbaasPostgresClusterToTable(item.ClusterResponse)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(item, clusterConverted,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Backup ID: %v", backupId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v from backup: %v", clusterId, backupId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	input := resources.CreateRestoreRequest{
		CreateRestoreRequest: psql.CreateRestoreRequest{
			BackupId: backupId,
		},
	}

	if viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime))))

		recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)))
		if err != nil {
			return err
		}

		input.SetRecoveryTargetTime(recoveryTargetTime)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Restoring Cluster from Backup..."))

	_, err := c.CloudApiDbaasPgsqlServices.Restores().Restore(clusterId, input)
	if err != nil {
		return err
	}
	if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("PostgreSQL Cluster successfully restored"))
	return nil
}

func RunClusterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		if err := ClusterDeleteAll(c); err != nil {
			return err
		}
		return nil
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v", clusterId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := c.CloudApiDbaasPgsqlServices.Clusters().Delete(clusterId)
	if err != nil {
		return err
	}
	if err = waitfor.WaitForDelete(c, waiter.ClusterDeleteInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
		return err
	}
	return nil
}

func ClusterDeleteAll(c *core.CommandConfig) error {

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))))
	}

	clusters, _, err := c.CloudApiDbaasPgsqlServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
	if err != nil {
		return err
	}

	dataOk, ok := clusters.GetItemsOk()
	if !ok || dataOk == nil {
		return fmt.Errorf("could not get items of Clusters")
	}

	if len(dataOk) <= 0 {
		return fmt.Errorf("no Clusters found")
	}

	var multiErr error
	for _, cluster := range dataOk {
		idOk, ok := cluster.GetIdOk()
		if !ok || idOk == nil {
			continue
		}

		clusterId := *idOk
		propertiesOk, ok := cluster.GetPropertiesOk()
		clusterName, ok := propertiesOk.GetDisplayNameOk()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v, name: %v", clusterId, *clusterName), viper.GetBool(constants.ArgForce)) {
			multiErr = errors.Join(multiErr, fmt.Errorf(confirm.UserDenied))
			continue
		}

		_, err = c.CloudApiDbaasPgsqlServices.Clusters().Delete(*idOk)
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *idOk, err))
			continue
		}

		if err = waitfor.WaitForDelete(c, waiter.ClusterDeleteInterrogator, *idOk); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *idOk, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func getCreateClusterRequest(c *core.CommandConfig) (*resources.CreateClusterRequest, error) {
	inputCluster := resources.CreateClusterRequest{}
	input := psql.CreateClusterProperties{}

	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("PostgresVersion: %v", pgsqlVersion))
	input.SetPostgresVersion(pgsqlVersion)

	syncMode := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagSyncMode)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("SynchronizationMode: %v", syncMode))
	input.SetSynchronizationMode(psql.SynchronizationMode(syncMode))

	replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Instances: %v", replicas))
	input.SetInstances(replicas)

	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Cores: %v", cpuCoreCount))
	input.SetCores(cpuCoreCount)

	// Convert Ram
	size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
	if err != nil {
		return nil, err
	}
	input.SetRam(int32(size))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Ram: %v[MB]", int32(size)))

	// Convert StorageSize
	storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
	if err != nil {
		return nil, err
	}
	input.SetStorageSize(int32(storageSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("StorageSize: %v[MB]", int32(storageSize)))
	storageType := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageType)))
	// "HDD" "SSD" "SSD Standard" "SSD Premium". "SSD" is deprecated and equivalent to "SSD Premium"
	if storageType == "SSD_PREMIUM" || storageType == "SSD PREMIUM" {
		storageType = "SSD Premium"
	}
	if storageType == "SSD_STANDARD" || storageType == "SSD STANDARD" {
		storageType = "SSD Standard"
	}
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("StorageType: %v", storageType))
	input.SetStorageType(psql.StorageType(storageType))

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Backup Location: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation))))
		input.SetBackupLocation(viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLocation)) {
		location := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLocation))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Location: %v", location))
		input.SetLocation(location)
	} else {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Location from VDC..."))
		vdc, _, err := c.CloudApiV6Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId)), cloudapiv6resources.QueryParams{})
		if err != nil {
			return nil, err
		}

		if properties, ok := vdc.GetPropertiesOk(); ok && properties != nil {
			if location, ok := properties.GetLocationOk(); ok && location != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Location: %v", *location))
				input.SetLocation(*location)
			}
		}
	}

	displayName := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DisplayName: %v", displayName))
	input.SetDisplayName(displayName)

	dbuser := psql.DBUser{}
	username := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDbUsername))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DBUser - Username: %v", username))
	dbuser.SetUsername(username)

	password := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDbPassword))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DBUser - Password: %v", password))
	dbuser.SetPassword(password)

	input.SetCredentials(dbuser)

	vdcConnection := psql.Connection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - DatacenterId: %v", vdcId))
	vdcConnection.SetDatacenterId(vdcId)

	lanId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLanId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - LanId: %v", lanId))
	vdcConnection.SetLanId(lanId)

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCidr)) {
		ip := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgCidr))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - Cidr: %v", ip))
		vdcConnection.SetCidr(ip)
	}

	input.SetConnections([]psql.Connection{vdcConnection})

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		maintenanceWindow := psql.MaintenanceWindow{}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceWindow - Time: %v", maintenanceTime))
			maintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay))
			maintenanceWindow.SetDayOfTheWeek(psql.DayOfTheWeek(maintenanceDay))
		}

		input.SetMaintenanceWindow(maintenanceWindow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) ||
		viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) {
		createRestoreRequest := psql.CreateRestoreRequest{}

		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)) {
			recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgRecoveryTime)))
			if err != nil {
				return nil, err
			}

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("From Backup - RecoveryTargetTime [RFC3339 format]: %v", recoveryTargetTime))
			createRestoreRequest.SetRecoveryTargetTime(recoveryTargetTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) {
			backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("From Backup - BackupId: %v", backupId))
			createRestoreRequest.SetBackupId(backupId)
		}

		input.SetFromBackup(createRestoreRequest)
	}

	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

func getPatchClusterRequest(c *core.CommandConfig) (*resources.PatchClusterRequest, error) {
	inputCluster := resources.PatchClusterRequest{}
	input := psql.PatchClusterProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Cores: %v", cpuCoreCount))
		input.SetCores(cpuCoreCount)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Instances: %v", replicas))
		input.SetInstances(replicas)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		// Convert Ram
		size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
		if err != nil {
			return nil, err
		}

		input.SetRam(int32(size))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Ram: %vMB", int32(size)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStorageSize)) {
		// Convert StorageSize
		storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
		if err != nil {
			return nil, err
		}

		input.SetStorageSize(int32(storageSize))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("StorageSize: %vMB", storageSize))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("PostgresVersion: %v", pgsqlVersion))
		input.SetPostgresVersion(pgsqlVersion)
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DisplayName: %v", displayName))
		input.SetDisplayName(displayName)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		maintenanceWindow := psql.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceTime: %v", maintenanceTime))
			maintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceDayOfWeek: %v", maintenanceDay))
			maintenanceWindow.SetDayOfTheWeek(psql.DayOfTheWeek(maintenanceDay))
		}

		input.SetMaintenanceWindow(maintenanceWindow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId)) || viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLanId)) || viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCidr)) {
		connection, err := getConnectionFromCluster(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(getConnectionMessage(connection)))

		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgDatacenterId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Datacenter Id: %v", lanId))
			connection.SetDatacenterId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgLanId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgLanId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Lan Id: %v", lanId))
			connection.SetLanId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, dbaaspg.ArgCidr)) {
			cidrId := viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgCidr))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Cidr: %v", cidrId))
			connection.SetCidr(cidrId)
		}

		input.SetConnections([]psql.Connection{connection})
	}

	if viper.GetBool(core.GetFlagName(c.NS, dbaaspg.ArgRemoveConnection)) {
		connection, err := getConnectionFromCluster(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Removing Connection with: %v...", getConnectionMessage(connection)))
		input.SetConnections([]psql.Connection{})

	}

	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

func getConnectionFromCluster(c *core.CommandConfig, clusterId string) (psql.Connection, error) {
	if c != nil {
		oldCluster, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(clusterId)
		if err != nil {
			return psql.Connection{}, err
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting properties from cluster with Id: %v", clusterId))
		if propertiesOk, ok := oldCluster.GetPropertiesOk(); ok && propertiesOk != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting connection.."))

			if connectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && connectionsOk != nil {
				for _, connectionOk := range connectionsOk {
					return connectionOk, nil
				}
			} else {
				return psql.Connection{}, errors.New("no connections found")
			}
		}
	}

	return psql.Connection{}, nil
}

func getConnectionMessage(connection psql.Connection) string {
	var msg string

	if datacenterOk, ok := connection.GetDatacenterIdOk(); ok && datacenterOk != nil {
		msg = fmt.Sprintf("DatacenterId: %v", *datacenterOk)
	}

	if lanOk, ok := connection.GetLanIdOk(); ok && lanOk != nil {
		msg = fmt.Sprintf("%v, LanId: %v", msg, *lanOk)
	}

	if cidrOk, ok := connection.GetCidrOk(); ok && cidrOk != nil {
		msg = fmt.Sprintf("%v, Cidr: %v", msg, *cidrOk)
	}

	return msg
}

// Output Printing

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "Location", "DatacenterId", "LanId", "Cidr", "Instances", "State"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "Location", "State", "PostgresVersion", "Instances", "Ram", "Cores",
		"StorageSize", "StorageType", "DatacenterId", "LanId", "Cidr", "MaintenanceWindow", "SynchronizationMode", "BackupLocation"}
)
