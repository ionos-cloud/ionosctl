package cluster

import (
	"context"
	"fmt"
	"math"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterUpdateCmd() *core.Command {
	ctx := context.TODO()

	update := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a PostgreSQL Cluster",
		LongDesc: `Use this command to update attributes of a PostgreSQL Cluster.

Required values to run command:

* Cluster Id
* DB Password`,
		Example: "ionosctl dbaas postgres-v2 cluster update --cluster-id <cluster-id> --db-password <password> --cores 4 --ram 8GB",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagDbPassword)
		},
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	update.AddStringFlag(constants.FlagDbPassword, constants.FlagDbPasswordShortPsql, "", "Password for the initial postgres user. Required because the API does not return it on GET requests", core.RequiredFlagOption())
	update.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "", "The PostgreSQL version of your cluster")

	update.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster. It has to be in the same location as the current datacenter")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagLanId, constants.FlagLanIdShortPsql, "", "The unique ID of the LAN to connect your cluster to")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(update.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagCidr, constants.FlagCidrShortPsql, "", "The IP and subnet for the cluster. Note the following unavailable IP range: 10.208.0.0/12. e.g.: 192.168.1.100/24")

	update.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 0, "The number of instances in your cluster. Minimum: 1, Maximum: 5")
	update.AddIntFlag(constants.FlagCores, "", 0, "The number of CPU cores per instance. Minimum: 1, Maximum: 62")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory per instance in GB. Minimum: 4, Maximum: 240. e.g. --ram 4, --ram 4GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB", "128GB", "240GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "The amount of storage per instance in GB. Minimum: 10, Maximum: 4096. e.g.: --storage-size 20, --storage-size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "500GB", "1TB", "2TB", "4TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The friendly name of your cluster")
	update.AddSetFlag(constants.FlagSyncModeV2, constants.FlagSyncModeShort, "", []string{"ASYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, "Replication mode")
	update.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")
	update.AddSetFlag(constants.FlagConnectionPooler, "", "", []string{"DISABLED", "TRANSACTION", "SESSION"}, "Connection pooling mode")
	update.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Enable collection and reporting of logs for this cluster")
	update.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Enable collection and reporting of metrics for this cluster")
	update.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "",
		"Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Must be specified together with --maintenance-day")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddSetFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "",
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		"Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time")
	update.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	update.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, table.ColsMessage(clusterCols))
	_ = update.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(clusterCols), cobra.ShellCompDirectiveNoFileComp
	})

	return update
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	c.Verbose(constants.ClusterId, clusterId)

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	newCluster, err := updateClusterProperties(c, clusterRead.Properties)
	if err != nil {
		return err
	}

	// Update (Ensure) Cluster
	c.Verbose("Updating Cluster...")
	clusterEnsure := psqlv2.NewClusterEnsure(clusterId, newCluster)

	item, _, err := client.Must().PostgresClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, clusterId); err != nil {
			return err
		}

		if item, _, err = client.Must().PostgresClientV2.ClustersApi.
			ClustersFindById(context.Background(), clusterId).Execute(); err != nil {
			return err
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(clusterCols, item, cols))
}

func updateClusterProperties(c *core.CommandConfig, input psqlv2.Cluster) (psqlv2.Cluster, error) {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		c.Verbose("Cores: %v", cpuCoreCount)
		input.Instances.Cores = cpuCoreCount
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
		c.Verbose("Instances: %v", replicas)
		input.Instances.Count = replicas
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		size, ok := convbytes.StrToUnitOk(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), convbytes.GB)
		if !ok {
			return input, fmt.Errorf("invalid value for Ram: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)))
		}

		if size < 0 || size > math.MaxInt32 {
			return input, fmt.Errorf("--ram value %vGB exceeds accepted int32 range: 0 - %d", size, math.MaxInt32)
		}
		input.Instances.Ram = int32(size)
		c.Verbose("Ram: %vGB", int32(size))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStorageSize)) {
		storageSize, ok := convbytes.StrToUnitOk(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), convbytes.GB)
		if !ok {
			return input, fmt.Errorf("invalid value for StorageSize: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)))
		}

		if storageSize < 0 || storageSize > math.MaxInt32 {
			return input, fmt.Errorf("--storage-size value %vGB exceeds accepted int32 range: 0 - %d", storageSize, math.MaxInt32)
		}
		input.Instances.StorageSize = int32(storageSize)
		c.Verbose("StorageSize: %vGB", storageSize)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		c.Verbose("PostgresVersion: %v", pgsqlVersion)
		input.SetVersion(pgsqlVersion)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
		c.Verbose("DisplayName: %v", displayName)
		input.SetName(displayName)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSyncModeV2)) {
		syncMode := viper.GetString(core.GetFlagName(c.NS, constants.FlagSyncModeV2))
		c.Verbose("ReplicationMode: %v", syncMode)
		input.SetReplicationMode(psqlv2.PostgresClusterReplicationMode(syncMode))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			c.Verbose("MaintenanceTime: %v", maintenanceTime)
			input.MaintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			c.Verbose("MaintenanceDayOfWeek: %v", maintenanceDay)
			input.MaintenanceWindow.SetDayOfTheWeek(psqlv2.DayOfTheWeek(maintenanceDay))
		}
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
		desc := viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))
		c.Verbose("Description: %v", desc)
		input.SetDescription(desc)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagConnectionPooler)) {
		cp := viper.GetString(core.GetFlagName(c.NS, constants.FlagConnectionPooler))
		c.Verbose("ConnectionPooler: %v", cp)
		input.SetConnectionPooler(cp)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLogsEnabled)) {
		logsEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagLogsEnabled))
		c.Verbose("LogsEnabled: %v", logsEnabled)
		input.SetLogsEnabled(logsEnabled)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMetricsEnabled)) {
		metricsEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagMetricsEnabled))
		c.Verbose("MetricsEnabled: %v", metricsEnabled)
		input.SetMetricsEnabled(metricsEnabled)
	}

	// Update Connection
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) {
			dcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
			c.Verbose("Updated Datacenter Id: %v", dcId)
			input.Connection.SetDatacenterId(dcId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
			c.Verbose("Updated Lan Id: %v", lanId)
			input.Connection.SetLanId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
			cidr := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
			c.Verbose("Updated Cidr: %v", cidr)
			input.Connection.SetPrimaryInstanceAddress(cidr)
		}
	}

	// Credentials - password is required on PUT because the API does not return it on GET
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDbPassword)) {
		password := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbPassword))
		credentials := input.GetCredentials()
		credentials.SetPassword(password)
		input.SetCredentials(credentials)
	}

	return input, nil
}
