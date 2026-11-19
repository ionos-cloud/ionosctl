package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCreateCmd() *core.Command {
	ctx := context.TODO()

	// Default maintenance window: a weekday between 10:00 and 16:00, picked from the current time.
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	seed := time.Now().UnixNano()
	hour := 10 + int(seed%7)
	defaultMaintenanceDay := workingDaysOfWeek[int(seed/7)%len(workingDaysOfWeek)]
	defaultMaintenanceTime := fmt.Sprintf("%02d:00:00", hour)

	create := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mongo-v2",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a MongoDB Cluster",
		LongDesc: `Use this command to create a new MongoDB Cluster with one primary and n-1 secondary instances. The --edition decides how the instances are sized:
  1. ENTERPRISE: pass --cores, --ram and --storage-size (defaults: 1 core, 8GB RAM, 10GB storage). Optionally enable sharding with --shards together with more than one replica per shard via --instances (sharding is permanent: shards can be added later, but never removed, and sharding cannot be disabled) and the MongoDB BI Connector with --bi-connector-enabled.
  2. BUSINESS: pass --template-id, the ID of an instance template that bundles cores, RAM and storage (see ` + "`ionosctl dbaas mongo-v2 template list`" + `). The template cannot be changed later.

Both editions require the connection and credential flags (--datacenter-id, --lan-id, --cidr, --user, --password, --database) plus --location. --cidr takes one IP address with netmask per instance; with sharding it takes exactly three, one per config server, whatever the instance and shard counts. The snapshot location (--backup-location) cannot be changed later either.

To create the cluster from an existing snapshot instead of empty, additionally pass --snapshot-id. The cluster version is taken from the snapshot (so --version is not needed; if given, it must match the snapshot's version). The snapshot is applied in full, up to the end of its recovery window, unless --recovery-time picks an earlier point within that window. Each instance needs at least the snapshot's RequiredStorageSize of storage: that value is in MiB, while --storage-size and template storage are in GB (1GB = 1024MiB, round up).`,
		Example: `# Create an empty ENTERPRISE cluster
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database>

# Create a sharded ENTERPRISE cluster: 2 shards of 3 replicas each, plus 3 addresses for the config servers
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --shards 2 --instances 3 --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr1>,<cidr2>,<cidr3> --user <username> --password <password> --database <database>

# Create a BUSINESS cluster with one instance, sized by an instance template
ionosctl dbaas mongo-v2 cluster create --location <location> --edition BUSINESS --template-id <template-id> --instances 1 --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database>

# Create a cluster from an existing snapshot, restored to a point inside its recovery window (version is taken from the snapshot)
ionosctl dbaas mongo-v2 cluster create --location <location> --edition ENTERPRISE --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --snapshot-id <snapshot-id> --recovery-time 2025-01-02T15:00:00Z`,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddSetFlag(constants.FlagEdition, "", "", editions,
		"The edition of the cluster, which cannot be changed later. ENTERPRISE is sized via --cores, --ram and --storage-size and supports sharding and the BI Connector; BUSINESS is sized by an instance template via --template-id",
		core.RequiredFlagOption(),
	)
	create.AddStringFlag(constants.FlagVersion, "", "7.0", "The MongoDB version of your cluster. Defaults to the snapshot's version when --snapshot-id is set",
		core.WithCompletion(completer.Versions, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	create.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances in the cluster (one primary and n-1 secondary), between 1 and 9. With sharding, this is the number of replicas per shard, and must be more than 1")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7", "9"}, cobra.ShellCompDirectiveNoFileComp
	})

	// ENTERPRISE sizing, sharding and BI Connector
	create.AddInt32Flag(constants.FlagCores, "", 1, "ENTERPRISE only: the number of CPU cores per instance, between 1 and 62")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "16"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRam, "", "8GB", "ENTERPRISE only: the amount of memory per instance, between 6GB and 240GB. e.g. --ram 8, --ram 8GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"8GB", "16GB", "32GB", "64GB", "128GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "10GB", "ENTERPRISE only: the size of the storage per instance, between 10GB and 4096GB. e.g. --storage-size 10 or --storage-size 10GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(constants.FlagShards, "", 0, "ENTERPRISE only: enable sharding and split the data across this many shards, between 2 and 5. --instances then sets the replicas per shard and must be more than 1, and --cidr takes exactly three addresses, one per config server. Sharding is permanent: shards can be added later but never removed, and it cannot be disabled")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagShards, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2", "3", "4", "5"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagBiConnectorEnabled, "", false, "ENTERPRISE only: enable the MongoDB BI Connector, reachable over the MySQL wire protocol at bi-connect.<mongo-host>:27015")

	// BUSINESS sizing
	create.AddUUIDFlag(constants.FlagTemplateId, "", "", "BUSINESS only, and required there: the ID of the instance template that sets the cores, RAM and storage of each instance. Cannot be changed later. See 'ionosctl dbaas mongo-v2 template list'",
		core.WithCompletion(func() []string { return completer.TemplateIds(mongo.MONGOEDITION_BUSINESS) }, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)

	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")

	// Connection
	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, "", "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringSliceFlag(constants.FlagCidr, "", nil, "The IP addresses with netmask of the cluster instances: one per instance, or exactly three with sharding. e.g.: 192.168.1.100/24,192.168.1.101/24,192.168.1.102/24. The range 10.208.0.0/12 is not available", core.RequiredFlagOption())

	// Snapshot configuration
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "eu-central-3", "The Object Storage location where snapshots will be stored. Cannot be changed later. For added data safety, use a different location than the cluster",
		core.WithCompletion(completer.SnapshotLocations, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	create.AddInt32Flag(constants.FlagRetentionDays, "", 7, "How many days snapshots are retained before being automatically deleted. Minimum: 1, Maximum: 365")

	create.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Forward this cluster's logs to the IONOS Logging Service. Has no effect unless that service is activated on your contract, which incurs additional costs")
	create.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Forward this cluster's metrics to the IONOS Monitoring Service. Has no effect unless that service is activated on your contract, which incurs additional costs")

	create.AddStringFlag(constants.FlagMaintenanceTime, "", defaultMaintenanceTime,
		"Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagMaintenanceDay, "", defaultMaintenanceDay,
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		"Day of the week for the MaintenanceWindow. Defaults to a random day during Mon-Fri")

	// Credentials
	create.AddStringFlag(constants.ArgUser, "", "", "The username of the initial MongoDB user: 2-16 letters, digits or underscores", core.RequiredFlagOption())
	create.AddStringFlag(constants.ArgPassword, "", "", "The password of the initial MongoDB user: 8-256 characters", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDatabase, "", "", "The name of the initial database, owned by the initial user: 2-63 letters, digits or underscores", core.RequiredFlagOption())

	// Restore from snapshot (optional)
	create.AddStringFlag(constants.FlagSnapshotId, "", "", "Create the cluster from this snapshot instead of empty. The snapshot must come from a cluster in the same location, and each instance needs at least its RequiredStorageSize (MiB) of storage. The connection/credential flags and --location are still required; the cluster version is taken from the snapshot",
		core.WithCompletion(completer.SnapshotIds, constants.MongoDBApiRegionalURL, constants.MongoDBLocations),
	)
	create.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "Advanced: with --snapshot-id, restore to a specific point in time WITHIN the snapshot's recovery window (PITR): 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). Defaults to the end of the window, i.e. the snapshot applied in full")

	create.Command.Flags().SortFlags = false

	return create
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	// The second set only adds --snapshot-id, so the usage hint also shows creating from a snapshot.
	base := []string{constants.FlagEdition, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr, constants.ArgUser, constants.ArgPassword, constants.FlagDatabase}
	if err := c.CheckRequiredFlagsSetsAndLocation(
		base,
		append(append([]string{}, base...), constants.FlagSnapshotId),
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) && !viper.IsSet(core.GetFlagName(c.NS, constants.FlagSnapshotId)) {
		return fmt.Errorf("--recovery-time requires --snapshot-id to be set")
	}
	if err := validateEditionFlags(c); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagShards)) && viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)) < 2 {
		return fmt.Errorf("--%s needs more than one replica per shard: set --%s (e.g. 3)",
			constants.FlagShards, constants.FlagInstances)
	}
	if err := validateInstanceAddresses(
		viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagCidr)),
		viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)),
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagShards)),
	); err != nil {
		return err
	}
	return validateRetentionDays(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagRetentionDays)))
}

func validateEditionFlags(c *core.PreCommandConfig) error {
	templateSet := viper.IsSet(core.GetFlagName(c.NS, constants.FlagTemplateId))
	switch mongo.MongoEdition(viper.GetString(core.GetFlagName(c.NS, constants.FlagEdition))) {
	case mongo.MONGOEDITION_BUSINESS:
		if !templateSet {
			return fmt.Errorf("--%s is required for the BUSINESS edition; see `ionosctl dbaas mongo-v2 template list`", constants.FlagTemplateId)
		}
		if set := setFlags(c.NS, enterpriseOnlyFlags); len(set) > 0 {
			return fmt.Errorf("%s cannot be used with the BUSINESS edition: its instances are sized by --%s, and it supports neither sharding nor the BI Connector",
				strings.Join(set, ", "), constants.FlagTemplateId)
		}
	case mongo.MONGOEDITION_ENTERPRISE:
		if templateSet {
			return fmt.Errorf("--%s is only valid for the BUSINESS edition; size ENTERPRISE clusters with --%s, --%s and --%s",
				constants.FlagTemplateId, constants.FlagCores, constants.FlagRam, constants.FlagStorageSize)
		}
	}
	return nil
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating Cluster...")

	cluster, _, err := client.Must().MongoClientV2.ClustersApi.ClustersPost(context.Background()).ClusterCreate(input).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}

func getCreateClusterRequest(c *core.CommandConfig) (mongo.ClusterCreate, error) {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))
	versionSet := c.Command.Command.Flags().Changed(constants.FlagVersion)
	recoveryTimeSet := viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime))

	// The snapshot supplies the default --version and --recovery-time.
	var snapshot mongo.SnapshotRead
	if snapshotId != "" && (!versionSet || !recoveryTimeSet) {
		var err error
		snapshot, _, err = client.Must().MongoClientV2.SnapshotsApi.SnapshotsFindById(context.Background(), snapshotId).Execute()
		if err != nil {
			return mongo.ClusterCreate{}, fmt.Errorf("could not read snapshot %s: %w", snapshotId, err)
		}
	}

	version := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	if snapshotId != "" && !versionSet {
		if v := snapshot.Properties.MongoVersion; v != nil && *v != "" {
			version = *v
		}
	}
	c.Verbose("Version: %v", version)

	name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	c.Verbose("Name: %v", name)

	instances := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))

	// Connection
	connection := mongo.MongoClusterConnection{
		DatacenterId:      viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId)),
		LanId:             viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId)),
		InstanceAddresses: viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagCidr)),
	}
	c.Verbose("Connection - DatacenterId: %v, LanId: %v, Cidr: %v", connection.DatacenterId, connection.LanId, connection.InstanceAddresses)

	// Snapshot configuration
	snapshotConfig := mongo.ClusterSnapshot{
		Location:      viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)),
		RetentionDays: viper.GetInt32(core.GetFlagName(c.NS, constants.FlagRetentionDays)),
	}
	c.Verbose("Snapshot - Location: %v, RetentionDays: %v", snapshotConfig.Location, snapshotConfig.RetentionDays)

	// Maintenance window (required) - always set from flags (which have defaults)
	maintenanceWindow := mongo.MaintenanceWindow{
		Time:         viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)),
		DayOfTheWeek: mongo.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))),
	}
	c.Verbose("MaintenanceWindow - Time: %v, Day: %v", maintenanceWindow.Time, maintenanceWindow.DayOfTheWeek)

	// Credentials
	password := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))
	credentials := mongo.MongoUser{
		Username: viper.GetString(core.GetFlagName(c.NS, constants.ArgUser)),
		Password: &password,
		Database: viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase)),
	}
	c.Verbose("Credentials - Username: %v, Database: %v", credentials.Username, credentials.Database)

	var props clusterProps
	var properties mongo.ClusterCreateProperties
	switch mongo.MongoEdition(viper.GetString(core.GetFlagName(c.NS, constants.FlagEdition))) {
	case mongo.MONGOEDITION_BUSINESS:
		templateId := viper.GetString(core.GetFlagName(c.NS, constants.FlagTemplateId))
		c.Verbose("Edition: BUSINESS, TemplateId: %v, Instances: %v", templateId, instances)
		business := mongo.NewBusinessClusterCreateWithDefaults()
		business.SetEdition(mongo.MONGOEDITION_BUSINESS)
		business.SetInstanceTemplateId(templateId)
		business.SetInstanceCount(instances)
		props, properties = business, mongo.BusinessClusterCreateAsClusterCreateProperties(business)
	default:
		enterprise, err := newEnterpriseClusterCreate(c, instances)
		if err != nil {
			return mongo.ClusterCreate{}, err
		}
		props, properties = enterprise, mongo.EnterpriseClusterCreateAsClusterCreateProperties(enterprise)
	}
	props.SetName(name)
	props.SetVersion(version)
	props.SetConnection(connection)
	props.SetSnapshot(snapshotConfig)
	props.SetMaintenanceWindow(maintenanceWindow)
	props.SetCredentials(credentials)
	applyOptionalFlags(c, props)

	// Restore from snapshot (optional)
	if snapshotId != "" {
		var recoveryTime time.Time
		if recoveryTimeSet {
			t, err := parseRecoveryTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return mongo.ClusterCreate{}, err
			}
			recoveryTime = t
		} else {
			recoveryTime = recoveryWindowEnd(snapshot.Properties)
		}
		restore := mongo.NewMongoRestoreFromSnapshot()
		restore.SourceSnapshotId = &snapshotId
		restore.RecoveryTargetDatetime = &mongo.IonosTime{Time: recoveryTime}
		c.Verbose("Restoring from snapshot %v to %v", snapshotId, recoveryTime.Format(time.RFC3339))
		props.SetRestoreFromSnapshot(*restore)
	}

	return mongo.ClusterCreate{Properties: properties}, nil
}

func newEnterpriseClusterCreate(c *core.CommandConfig, instances int32) (*mongo.EnterpriseClusterCreate, error) {
	ram, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), constants.FlagRam)
	if err != nil {
		return nil, err
	}
	storage, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), constants.FlagStorageSize)
	if err != nil {
		return nil, err
	}
	instanceConfig := mongo.InstanceConfiguration{
		Count:       instances,
		Cores:       viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores)),
		Ram:         ram,
		StorageSize: storage,
	}
	c.Verbose("Edition: ENTERPRISE, Instances: %v, Cores: %v, Ram: %vGB, StorageSize: %vGB", instanceConfig.Count, instanceConfig.Cores, ram, storage)

	enterprise := mongo.NewEnterpriseClusterCreateWithDefaults()
	enterprise.SetEdition(mongo.MONGOEDITION_ENTERPRISE)
	enterprise.SetInstances(instanceConfig)
	applyEnterpriseFeatureFlags(c, enterprise)
	return enterprise, nil
}
