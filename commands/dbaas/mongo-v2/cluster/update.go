package cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterUpdateCmd() *core.Command {
	ctx := context.TODO()

	update := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mongo-v2",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a MongoDB Cluster",
		LongDesc: `Use this command to update attributes of a MongoDB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

The edition, the instance template of a BUSINESS cluster, the snapshot location, the datacenter and the LAN are fixed at creation. A BUSINESS cluster can change its --instances; an ENTERPRISE cluster its --instances, --cores, --ram and --storage-size (storage can only grow). On an ENTERPRISE cluster created with sharding, --shards adds shards: shards cannot be removed, and sharding cannot be turned on or off after creation. --bi-connector-enabled toggles the BI Connector of an ENTERPRISE cluster.

When the number of instances changes, pass --cidr in the same command so it keeps listing one address per instance (a sharded cluster always has exactly three).

The version can only be upgraded, to one of the CanUpgradeTo versions of ` + "`ionosctl dbaas mongo-v2 version list`" + `. Upgrading asks for confirmation: snapshots taken before the upgrade can no longer be restored in place.

The initial user's credentials stay unchanged unless --password is set. The API never returns the password, so changing --user or --database requires --password too.

Required values to run command:

* Cluster Id`,
		Example: `ionosctl dbaas mongo-v2 cluster update --location <location> --cluster-id <cluster-id> --cores 2 --ram 16GB --storage-size 20GB

# Scale to three instances, one address per instance
ionosctl dbaas mongo-v2 cluster update --location <location> --cluster-id <cluster-id> --instances 3 --cidr <cidr1>,<cidr2>,<cidr3>`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			if err := c.CheckRequiredFlagsAndLocation(constants.FlagClusterId); err != nil {
				return err
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRetentionDays); viper.IsSet(fn) {
				if err := validateRetentionDays(viper.GetInt32(fn)); err != nil {
					return err
				}
			}
			return validateCredentialFlags(c)
		},
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")
	update.AddStringFlag(constants.FlagVersion, "", "", "The MongoDB version to upgrade to: one of the CanUpgradeTo versions of the current one. Downgrades are not supported. Asks for confirmation, as snapshots taken before the upgrade can no longer be restored in place",
		core.WithCompletion(completer.Versions, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	update.AddInt32Flag(constants.FlagInstances, "", 0, "The total number of instances in the cluster (one primary and n-1 secondary), between 1 and 9. With sharding, this is the number of replicas per shard")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7", "9"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(constants.FlagCores, "", 0, "ENTERPRISE only: the number of CPU cores per instance, between 1 and 62")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "16"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagRam, "", "", "ENTERPRISE only: the amount of memory per instance, between 6GB and 240GB. e.g. --ram 8, --ram 8GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"8GB", "16GB", "32GB", "64GB", "128GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "ENTERPRISE only: the size of the storage per instance, up to 4096GB. Can only be increased. e.g. --storage-size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(constants.FlagShards, "", 0, "ENTERPRISE only: the new number of shards, up to 5, for a cluster created with sharding. Shards can only be added: they cannot be removed, and sharding cannot be enabled on an existing cluster")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagShards, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"3", "4", "5"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.FlagBiConnectorEnabled, "", false, "ENTERPRISE only: enable or disable the MongoDB BI Connector")
	update.AddStringSliceFlag(constants.FlagCidr, "", nil, "The IP addresses with netmask of the cluster instances, replacing the current ones: one per instance, or exactly three with sharding")
	update.AddInt32Flag(constants.FlagRetentionDays, "", 0, "How many days snapshots are retained. Minimum: 1, Maximum: 365. Reducing it purges the existing snapshots outside the new retention window")
	update.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Forward this cluster's logs to the IONOS Logging Service. Has no effect unless that service is activated on your contract, which incurs additional costs")
	update.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Forward this cluster's metrics to the IONOS Monitoring Service. Has no effect unless that service is activated on your contract, which incurs additional costs")
	update.AddStringFlag(constants.FlagMaintenanceTime, "", "",
		"Time for the MaintenanceWindow. e.g.: 16:30:59. Must be specified together with --maintenance-day")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddSetFlag(constants.FlagMaintenanceDay, "", "",
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		"Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time")

	update.AddStringFlag(constants.ArgPassword, "", "", "New password for the initial MongoDB user: 8-256 characters. When omitted, the user's credentials are left unchanged")
	update.AddStringFlag(constants.ArgUser, "", "", "Username for the initial MongoDB user. Requires --password. Defaults to the cluster's current username")
	update.AddStringFlag(constants.FlagDatabase, "", "", "Database for the initial MongoDB user. Requires --password. Defaults to the cluster's current database")

	return update
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	c.Verbose(constants.ClusterId, clusterId)

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().MongoClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	cluster := clusterRead.Properties
	props, err := editionProps(&cluster)
	if err != nil {
		return err
	}

	if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
		newVersion, currentVersion := viper.GetString(fn), props.GetVersion()
		if newVersion != "" && newVersion != currentVersion {
			prompt := fmt.Sprintf(
				"changing version %s -> %s: snapshots taken on %s can no longer be restored in place onto this cluster "+
					"(the API requires the snapshot and cluster versions to match). To recover pre-upgrade data later, create a NEW "+
					"cluster from one of those snapshots (`cluster create --snapshot-id <id>`, which takes version %s from the snapshot) and then upgrade it. Continue",
				currentVersion, newVersion, currentVersion, currentVersion)
			if !confirm.FAsk(c.Command.Command.InOrStdin(), prompt, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}
		}
	}

	if err := updateClusterProperties(c, &cluster); err != nil {
		return err
	}

	// Update (Ensure) Cluster
	c.Verbose("Updating Cluster...")
	clusterEnsure := mongo.NewClusterEnsure(clusterId, cluster)

	item, _, err := client.Must().MongoClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(clusterCols, item, cols))
}

func updateClusterProperties(c *core.CommandConfig, cluster *mongo.Cluster) error {
	props, err := editionProps(cluster)
	if err != nil {
		return err
	}
	currentAddresses := props.GetConnection().InstanceAddresses

	var instances int32
	var sharded bool
	switch {
	case cluster.EnterpriseCluster != nil:
		if err := updateEnterpriseProperties(c, cluster.EnterpriseCluster); err != nil {
			return err
		}
		instances = cluster.EnterpriseCluster.Instances.Count
		sharding := cluster.EnterpriseCluster.GetSharding()
		sharded = sharding.GetEnabled()
	case cluster.BusinessCluster != nil:
		if set := setFlags(c.NS, enterpriseOnlyFlags); len(set) > 0 {
			return fmt.Errorf("%s cannot be used on a BUSINESS cluster: its instances are sized by its instance template, which is fixed at creation, and it supports neither sharding nor the BI Connector",
				strings.Join(set, ", "))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.IsSet(fn) {
			count := viper.GetInt32(fn)
			c.Verbose("Instances: %v", count)
			cluster.BusinessCluster.SetInstanceCount(count)
		}
		instances = cluster.BusinessCluster.GetInstanceCount()
	}

	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		name := viper.GetString(fn)
		c.Verbose("Name: %v", name)
		props.SetName(name)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
		version := viper.GetString(fn)
		c.Verbose("Version: %v", version)
		props.SetVersion(version)
	}
	applyOptionalFlags(c, props)

	connection := props.GetConnection()
	if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
		connection.InstanceAddresses = viper.GetStringSlice(fn)
		c.Verbose("Cidr: %v", connection.InstanceAddresses)
		props.SetConnection(connection)
	}
	if len(setFlags(c.NS, []string{constants.FlagInstances, constants.FlagCidr, constants.FlagShards})) > 0 {
		if err := validateInstanceAddresses(connection.InstanceAddresses, instances, sharded); err != nil {
			return fmt.Errorf("%w; the cluster currently has --%s %s", err, constants.FlagCidr, strings.Join(currentAddresses, ","))
		}
	}

	if fn := core.GetFlagName(c.NS, constants.FlagRetentionDays); viper.IsSet(fn) {
		snapshot := props.GetSnapshot()
		snapshot.RetentionDays = viper.GetInt32(fn)
		c.Verbose("RetentionDays: %v", snapshot.RetentionDays)
		props.SetSnapshot(snapshot)
	}

	maintenanceWindow := props.GetMaintenanceWindow()
	if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.IsSet(fn) {
		maintenanceWindow.Time = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.IsSet(fn) {
		maintenanceWindow.DayOfTheWeek = mongo.DayOfTheWeek(viper.GetString(fn))
	}
	props.SetMaintenanceWindow(maintenanceWindow)

	return applyCredentialsFromFlags(c, cluster)
}

func updateEnterpriseProperties(c *core.CommandConfig, cluster *mongo.EnterpriseCluster) error {
	if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.IsSet(fn) {
		instances := viper.GetInt32(fn)
		c.Verbose("Instances: %v", instances)
		cluster.Instances.Count = instances
	}
	if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(fn) {
		cores := viper.GetInt32(fn)
		c.Verbose("Cores: %v", cores)
		cluster.Instances.Cores = cores
	}
	if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.IsSet(fn) {
		ram, err := toInt32GB(viper.GetString(fn), constants.FlagRam)
		if err != nil {
			return err
		}
		c.Verbose("Ram: %vGB", ram)
		cluster.Instances.Ram = ram
	}
	if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.IsSet(fn) {
		storage, err := toInt32GB(viper.GetString(fn), constants.FlagStorageSize)
		if err != nil {
			return err
		}
		c.Verbose("StorageSize: %vGB", storage)
		cluster.Instances.StorageSize = storage
	}
	if fn := core.GetFlagName(c.NS, constants.FlagShards); viper.IsSet(fn) {
		if err := validateShardIncrease(cluster.GetSharding(), viper.GetInt32(fn)); err != nil {
			return err
		}
	}
	applyEnterpriseFeatureFlags(c, cluster)
	return nil
}

func validateShardIncrease(current mongo.MongoSharding, shards int32) error {
	if !current.GetEnabled() {
		return fmt.Errorf("--%s can only add shards to a cluster created with sharding: sharding cannot be enabled on an existing cluster", constants.FlagShards)
	}
	if cur, ok := current.GetShardsOk(); ok && shards < *cur {
		return fmt.Errorf("--%s cannot remove shards: the cluster has %d, requested %d", constants.FlagShards, *cur, shards)
	}
	return nil
}
