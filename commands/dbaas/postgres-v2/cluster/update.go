package cluster

import (
	"context"
	"fmt"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterUpdateCmd() *core.Command {
	ctx := context.TODO()

	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	update := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a PostgreSQL Cluster",
		LongDesc: `Use this command to update attributes of a PostgreSQL Cluster.

Required values to run command:

* Cluster Id`,
		Example: "ionosctl dbaas postgres-v2 cluster update --cluster-id <cluster-id> --cores 4 --ram 8GB",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
		},
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			clusters, err := Clusters()
			if err != nil {
				return []string{}
			}
			return functional.Map(clusters.Items, func(c psqlv2.ClusterRead) string {
				return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
					c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)
			})
		}, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
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

	update.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 0, "The number of instances in your cluster. Minimum: 1. Maximum: 5")
	update.AddIntFlag(constants.FlagCores, "", 0, "The number of CPU cores per instance")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory per instance in GB. e.g. --ram 4096, --ram 4096MB, --ram 4GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "The amount of storage per instance in GB. e.g.: --storage-size 20480 or --storage-size 20480MB or --storage-size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(constants.FlagSyncMode, constants.FlagSyncModeShort, "", "Replication mode: ASYNCHRONOUS, SYNCHRONOUS, STRICTLY_SYNCHRONOUS")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "",
		"Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Must be specified together with --maintenance-day")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "",
		"Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	update.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = update.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return update
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))

	// Fetch existing cluster
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Cluster..."))
	clusterRead, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	newCluster, err := updateClusterProperties(c, clusterRead.Properties)
	if err != nil {
		return err
	}

	// Update (Ensure) Cluster
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating Cluster..."))
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
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresV2Cluster, item,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func updateClusterProperties(c *core.CommandConfig, input psqlv2.Cluster) (psqlv2.Cluster, error) {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Cores: %v", cpuCoreCount))
		input.Instances.Cores = cpuCoreCount
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Instances: %v", replicas))
		input.Instances.Count = replicas
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
		if err != nil {
			return input, err
		}

		input.Instances.Ram = int32(size)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Ram: %vMB", int32(size)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStorageSize)) {
		storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
		if err != nil {
			return input, err
		}

		input.Instances.StorageSize = int32(storageSize)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("StorageSize: %vMB", storageSize))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("PostgresVersion: %v", pgsqlVersion))
		input.SetVersion(pgsqlVersion)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DisplayName: %v", displayName))
		input.SetName(displayName)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSyncMode)) {
		syncMode := viper.GetString(core.GetFlagName(c.NS, constants.FlagSyncMode))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("ReplicationMode: %v", syncMode))
		input.SetReplicationMode(psqlv2.PostgresClusterReplicationMode(syncMode))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceTime: %v", maintenanceTime))
			input.MaintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceDayOfWeek: %v", maintenanceDay))
			input.MaintenanceWindow.SetDayOfTheWeek(psqlv2.DayOfTheWeek(maintenanceDay))
		}
	}

	// Update Connection
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) {
			dcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Datacenter Id: %v", dcId))
			input.Connection.SetDatacenterId(dcId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Lan Id: %v", lanId))
			input.Connection.SetLanId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
			cidr := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updated Cidr: %v", cidr))
			input.Connection.SetPrimaryInstanceAddress(cidr)
		}
	}

	return input, nil
}
