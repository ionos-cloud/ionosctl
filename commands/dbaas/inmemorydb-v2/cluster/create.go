package cluster

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const persistenceModeLong = `PersistenceMode:
None: Data is inMemory only and will not be persisted. Useful for cache only applications.
AOF (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset.
RDB: RDB persistence performs snapshots of the current in memory state.
RDB_AOF: Both RDB and AOF persistence are enabled.

EvictionPolicy:
noeviction: No eviction policy is used. In-Memory DB will never remove any data. If the memory limit is reached, an error will be returned on write operations.
allkeys-lru: The least recently used keys will be removed first.
allkeys-lfu: The least frequently used keys will be removed first.
allkeys-random: Random keys will be removed.
volatile-lru: The least recently used keys will be removed first, but only among keys with the expire field set to true.
volatile-lfu: The least frequently used keys will be removed first, but only among keys with the expire field set to true.
volatile-random: Random keys will be removed, but only among keys with the expire field set to true.
volatile-ttl: The key with the nearest time to live will be removed first, but only among keys with the expire field set to true.`

func ClusterCreateCmd() *core.Command {
	ctx := context.TODO()

	// Generate random maintenance window defaults (Mon-Fri, 10:00-16:00)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7)
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	defaultMaintenanceDay := workingDaysOfWeek[r.Intn(len(workingDaysOfWeek))]
	defaultMaintenanceTime := fmt.Sprintf("%02d:00:00", hour)

	create := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-inmemorydb-v2",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an In-Memory DB Cluster",
		LongDesc: `Use this command to create a new In-Memory DB Cluster. The mode is determined by the number of replicas: one replica is standalone, everything else is a replication in leader-follower mode with one active and n-1 passive replicas.

` + persistenceModeLong,
		Example:    "ionosctl dbaas in-memory-db-v2 cluster create --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --version <version>",
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, "", "8.0", "The In-Memory DB version of your Cluster", core.RequiredFlagOption(),
		core.WithCompletion(completer.Versions, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	create.AddIntFlag(constants.FlagReplicas, "", 1, "The total number of replicas in the cluster (one active and n-1 passive). In case of a standalone instance, the value is 1")
	create.AddIntFlag(constants.FlagCores, "", 1, "The number of CPU cores per instance")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "12", "16", "24", "31"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance in gigabytes (GB). e.g. --ram 4, --ram 4GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB", "128GB", "256GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagPersistenceMode, "", "RDB", []string{"None", "AOF", "RDB", "RDB_AOF"}, "Specifies how and if data is persisted (refer to the long description for more details)")
	create.AddSetFlag(constants.FlagEvictionPolicy, "", "allkeys-lru",
		[]string{"noeviction", "allkeys-lru", "allkeys-lfu", "allkeys-random", "volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl"}, "The eviction policy for the cluster (refer to the long description for more details)")
	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")

	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, "", "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCidr, "", "", "The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.210.0.0/16 10.212.0.0/14. e.g.: 192.168.1.100/24", core.RequiredFlagOption())

	// Snapshot configuration
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "eu-central-4", "The Object Storage location where snapshots (backups) will be stored. For added data safety, use a different location than the cluster",
		core.WithCompletion(completer.SnapshotLocations, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	create.AddInt32Flag(constants.FlagRetentionDays, "", 7, "The number of days snapshots are retained before being automatically deleted")
	create.AddIntSliceFlag(constants.FlagSnapshotHours, "", []int{4}, "Hours of the day (UTC, 0-23) at which snapshots are scheduled to be taken. At least one hour must be specified")

	create.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Enable collection and reporting of logs for this cluster")
	create.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Enable collection and reporting of metrics for this cluster")

	create.AddStringFlag(constants.FlagMaintenanceTime, "", defaultMaintenanceTime,
		"Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random time during 10:00-16:00")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddSetFlag(constants.FlagMaintenanceDay, "", defaultMaintenanceDay,
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		"Day of the week for the MaintenanceWindow. Defaults to a random day during Mon-Fri")

	// Credentials
	create.AddStringFlag(constants.ArgUser, "", "", "The initial username", core.RequiredFlagOption())
	create.AddStringFlag(constants.ArgPassword, "", "", "Password for the initial user. Plaintext is hashed (SHA-256) client-side before sending, as the API only accepts hashed passwords; a value that is already a SHA-256 hash is sent as-is", core.RequiredFlagOption())

	// Restore from snapshot (optional)
	create.AddStringFlag(constants.FlagSnapshotId, "", "", "If set, create the cluster restored from the specified snapshot",
		core.WithCompletion(completer.SnapshotIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	create.AddStringFlag(constants.FlagRecoveryTime, "", "", "Together with --snapshot-id, an ISO 8601 timestamp to restore from the most recent snapshot taken at or before that time")

	return create
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	if err := c.CheckRequiredFlagsAndLocation(constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
		constants.ArgUser, constants.ArgPassword); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) && !viper.IsSet(core.GetFlagName(c.NS, constants.FlagSnapshotId)) {
		return fmt.Errorf("--recovery-time requires --snapshot-id to be set")
	}
	return nil
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating Cluster...")

	cluster, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersPost(context.Background()).ClusterCreate(input).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}

func getCreateClusterRequest(c *core.CommandConfig) (inmemorydb.ClusterCreate, error) {
	inputCluster := inmemorydb.ClusterCreate{}
	input := inmemorydb.ClusterCreateProperties{}

	version := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	c.Verbose("Version: %v", version)
	input.Version = version

	displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	c.Verbose("Name: %v", displayName)
	input.Name = displayName

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
		desc := viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))
		c.Verbose("Description: %v", desc)
		input.SetDescription(desc)
	}

	instanceConfig := inmemorydb.InstanceConfiguration{}
	replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagReplicas))
	c.Verbose("Replicas: %v", replicas)
	instanceConfig.Count = replicas

	cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	c.Verbose("Cores: %v", cores)
	instanceConfig.Cores = cores

	ram, err := ramToInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)))
	if err != nil {
		return inputCluster, err
	}
	instanceConfig.Ram = ram
	c.Verbose("Ram: %vGB", ram)
	input.Instances = instanceConfig

	persistenceMode := inmemorydb.PersistenceMode(viper.GetString(core.GetFlagName(c.NS, constants.FlagPersistenceMode)))
	c.Verbose("PersistenceMode: %v", persistenceMode)
	input.PersistenceMode = &persistenceMode

	evictionPolicy := inmemorydb.EvictionPolicy(viper.GetString(core.GetFlagName(c.NS, constants.FlagEvictionPolicy)))
	c.Verbose("EvictionPolicy: %v", evictionPolicy)
	input.EvictionPolicy = evictionPolicy

	// Connection
	connection := inmemorydb.ClusterConnection{}
	connection.DatacenterId = viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
	connection.LanId = viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
	connection.PrimaryInstanceAddress = viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
	c.Verbose("Connection - DatacenterId: %v, LanId: %v, Cidr: %v", connection.DatacenterId, connection.LanId, connection.PrimaryInstanceAddress)
	input.Connection = connection

	// Snapshot configuration
	snapshotConfig := inmemorydb.SnapshotConfiguration{}
	snapshotConfig.Location = viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation))
	snapshotConfig.RetentionDays = viper.GetInt32(core.GetFlagName(c.NS, constants.FlagRetentionDays))
	snapshotConfig.SnapshotHours = intsToInt32s(viper.GetIntSlice(core.GetFlagName(c.NS, constants.FlagSnapshotHours)))
	c.Verbose("Snapshot - Location: %v, RetentionDays: %v, SnapshotHours: %v", snapshotConfig.Location, snapshotConfig.RetentionDays, snapshotConfig.SnapshotHours)
	input.Snapshot = snapshotConfig

	// Maintenance window (required) - always set from flags (which have defaults)
	maintenanceWindow := inmemorydb.MaintenanceWindow{}
	maintenanceWindow.Time = viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
	maintenanceWindow.DayOfTheWeek = inmemorydb.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)))
	c.Verbose("MaintenanceWindow - Time: %v, Day: %v", maintenanceWindow.Time, maintenanceWindow.DayOfTheWeek)
	input.MaintenanceWindow = maintenanceWindow

	// Credentials
	credentials := inmemorydb.ClusterCredentials{}
	credentials.Username = viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	credentials.Password = buildHashedPassword(viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword)))
	input.Credentials = credentials

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

	// Restore from snapshot (optional)
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSnapshotId)) {
		snapshotId := viper.GetString(core.GetFlagName(c.NS, constants.FlagSnapshotId))
		restore := inmemorydb.NewRestoreClusterFromSnapshot(snapshotId)
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
			t, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return inputCluster, fmt.Errorf("invalid recovery-time format (expected RFC3339, e.g. 2024-01-15T10:00:00Z): %w", err)
			}
			restore.RecoveryTargetDatetime = &inmemorydb.IonosTime{Time: t}
		}
		c.Verbose("Restoring from snapshot: %v", snapshotId)
		input.RestoreFromSnapshot = &inmemorydb.ClusterRestoreFromSnapshot{RestoreClusterFromSnapshot: restore}
	}

	inputCluster.Properties = input
	return inputCluster, nil
}

// intsToInt32s converts a []int (from viper) to []int32 for the SDK.
func intsToInt32s(in []int) []int32 {
	out := make([]int32, len(in))
	for i, v := range in {
		out[i] = int32(v)
	}
	return out
}

// ramToInt32GB converts a RAM flag value (e.g. "4GB") to an int32 amount of GB.
func ramToInt32GB(raw string) (int32, error) {
	size, ok := convbytes.StrToUnitOk(raw, convbytes.GB)
	if !ok {
		return 0, fmt.Errorf("invalid value for Ram: %v", raw)
	}
	if size < 0 || size > math.MaxInt32 {
		return 0, fmt.Errorf("--ram value %vGB exceeds accepted int32 range: 0 - %d", size, math.MaxInt32)
	}
	return int32(size), nil
}
