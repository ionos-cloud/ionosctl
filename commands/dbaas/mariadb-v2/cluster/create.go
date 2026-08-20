package cluster

import (
	"context"
	"fmt"
	"math"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCreateCmd() *core.Command {
	ctx := context.TODO()

	// Spread the default maintenance window across customers (Mon-Fri, 10:00-16:00)
	// using the current time as entropy. This is a cosmetic default, not a
	// security-sensitive value, so a time-derived spread is sufficient (and avoids
	// a non-crypto PRNG).
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	seed := time.Now().UnixNano()
	hour := 10 + int(seed%7)
	defaultMaintenanceDay := workingDaysOfWeek[int(seed/7)%len(workingDaysOfWeek)]
	defaultMaintenanceTime := fmt.Sprintf("%02d:00:00", hour)

	create := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb-v2",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a MariaDB Cluster",
		LongDesc: `Use this command to create a new MariaDB Cluster. The mode is determined by the number of instances: one instance is standalone, everything else is a replication with one primary and n-1 secondary instances.

There are two ways to create a cluster, both requiring the same connection and credential flags (--datacenter-id, --lan-id, --cidr, --user, --password, --database) plus --location:
  1. Empty cluster: pass --version (defaults to a supported version) and sizing flags (--instances, --cores, --ram, --storage-size).
  2. From a backup: additionally pass --backup-id. The cluster version is taken from the backup (so --version is not needed; if given, it must match the backup's version). Optionally pass --recovery-time to restore to a point in time within the backup's window.`,
		Example: `# Create an empty cluster
ionosctl dbaas mariadb-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --version <version>

# Create a cluster from an existing backup (version is taken from the backup)
ionosctl dbaas mariadb-v2 cluster create --location <location> --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --user <username> --password <password> --database <database> --backup-id <backup-id>`,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, "", "10.11", "The MariaDB version of your cluster. Ignored when --backup-id is set (the backup's version is used)", core.RequiredFlagOption(),
		core.WithCompletion(completer.Versions, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	create.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances in the cluster (one primary and n-1 secondary). For a standalone instance, use 1")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddInt32Flag(constants.FlagCores, "", 1, "The number of CPU cores per instance")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "16"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance. e.g. --ram 4, --ram 4GB. Minimum 4GB", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "10GB", "The size of the storage per instance. e.g. --storage-size 10 or --storage-size 10GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "The friendly name of your cluster")

	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, "", "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCidr, "", "", "The IP and subnet for the cluster. All IPs must be in a /24 network. e.g.: 192.168.1.100/24", core.RequiredFlagOption())

	// Backup configuration
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "eu-central-4", "The Object Storage location where backups will be stored. For added data safety, use a different location than the cluster",
		core.WithCompletion(completer.BackupLocations, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	create.AddInt32Flag(constants.FlagBackupRetentionDays, "", 30, "Configures how many days cluster backups are retained before being automatically deleted. Minimum: 1, Maximum: 365")

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
	create.AddStringFlag(constants.ArgPassword, "", "", "Password for the initial user", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDatabase, "", "", "The name of the initial database created for the user", core.RequiredFlagOption())

	// Restore from backup (optional)
	create.AddStringFlag(constants.FlagBackupId, "", "", "Create the cluster from this backup instead of empty. The connection/credential flags and --location are still required; the cluster version is taken from the backup",
		core.WithCompletion(completer.BackupIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	create.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "Advanced: with --backup-id, restore to a specific point in time WITHIN the backup's recovery window (PITR): 'now', a date, a date-time, or an RFC3339 timestamp (no timezone = UTC). Defaults to the latest point in the window")

	create.Command.Flags().SortFlags = false

	return create
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	// Two valid shapes, both requiring the connection/credential flags (+ --location,
	// prepended by the helper). The second additionally lists --backup-id so the
	// "create from backup" variant shows up in the usage hint.
	base := []string{constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr, constants.ArgUser, constants.ArgPassword, constants.FlagDatabase}
	if err := c.CheckRequiredFlagsSetsAndLocation(
		base,
		append(append([]string{}, base...), constants.FlagBackupId),
	); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) && !viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) {
		return fmt.Errorf("--recovery-time requires --backup-id to be set")
	}
	return validateBackupRetentionDays(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays)))
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating Cluster...")

	cluster, _, err := client.Must().MariaClientV2.ClustersApi.ClustersPost(context.Background()).ClusterCreate(input).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}

func getCreateClusterRequest(c *core.CommandConfig) (mariadb.ClusterCreate, error) {
	inputCluster := mariadb.ClusterCreate{}
	input := mariadb.ClusterCreateProperties{}

	version := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	// When creating from a backup, the cluster version must match the backup's
	// version. Derive it from the backup (unless the user set --version explicitly)
	// so the default --version does not cause a version-mismatch error.
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) && !c.Command.Command.Flags().Changed(constants.FlagVersion) {
		backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
		backup, _, err := client.Must().MariaClientV2.BackupsApi.BackupsFindById(context.Background(), backupId).Execute()
		if err != nil {
			return inputCluster, fmt.Errorf("could not read backup %s to derive the cluster version (pass --version explicitly): %w", backupId, err)
		}
		if v := backup.Properties.MariadbClusterVersion; v != nil && *v != "" {
			version = *v
		}
	}
	c.Verbose("Version: %v", version)
	input.Version = version

	name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	c.Verbose("Name: %v", name)
	input.Name = name

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
		desc := viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))
		c.Verbose("Description: %v", desc)
		input.SetDescription(desc)
	}

	instanceConfig := mariadb.InstanceConfiguration{}
	instanceConfig.Count = viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
	instanceConfig.Cores = viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))

	ram, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), "ram")
	if err != nil {
		return inputCluster, err
	}
	instanceConfig.Ram = ram

	storage, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), "storage-size")
	if err != nil {
		return inputCluster, err
	}
	instanceConfig.StorageSize = storage
	c.Verbose("Instances: %v, Cores: %v, Ram: %vGB, StorageSize: %vGB", instanceConfig.Count, instanceConfig.Cores, ram, storage)
	input.Instances = instanceConfig

	// Connection
	connection := mariadb.MariadbClusterConnection{}
	connection.DatacenterId = viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
	connection.LanId = viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
	connection.PrimaryInstanceAddress = viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
	c.Verbose("Connection - DatacenterId: %v, LanId: %v, Cidr: %v", connection.DatacenterId, connection.LanId, connection.PrimaryInstanceAddress)
	input.Connection = connection

	// Backup configuration
	backupConfig := mariadb.ClusterBackup{}
	backupConfig.Location = viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation))
	backupConfig.RetentionDays = viper.GetInt32(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays))
	c.Verbose("Backup - Location: %v, RetentionDays: %v", backupConfig.Location, backupConfig.RetentionDays)
	input.Backup = backupConfig

	// Maintenance window (required) - always set from flags (which have defaults)
	maintenanceWindow := mariadb.MaintenanceWindow{}
	maintenanceWindow.Time = viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
	maintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)))
	c.Verbose("MaintenanceWindow - Time: %v, Day: %v", maintenanceWindow.Time, maintenanceWindow.DayOfTheWeek)
	input.MaintenanceWindow = maintenanceWindow

	// Credentials
	credentials := mariadb.MariadbUser{}
	credentials.Username = viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	credentials.Password = viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))
	credentials.Database = viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))
	c.Verbose("Credentials - Username: %v, Database: %v", credentials.Username, credentials.Database)
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

	// Restore from backup (optional)
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) {
		backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
		restore := mariadb.NewMariadbRestoreClusterFromBackup()
		restore.SourceBackupId = &backupId
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
			t, err := parseRecoveryTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return inputCluster, err
			}
			restore.RecoveryTargetDatetime = &mariadb.IonosTime{Time: t}
		}
		c.Verbose("Restoring from backup: %v", backupId)
		input.RestoreFromBackup = restore
	}

	inputCluster.Properties = input
	return inputCluster, nil
}

// toInt32GB converts a size flag value (e.g. "4GB") to an int32 amount of GB.
func toInt32GB(raw, flag string) (int32, error) {
	size, ok := convbytes.StrToUnitOk(raw, convbytes.GB)
	if !ok {
		return 0, fmt.Errorf("invalid value for --%s: %v", flag, raw)
	}
	if size < 0 || size > math.MaxInt32 {
		return 0, fmt.Errorf("--%s value %vGB exceeds accepted int32 range: 0 - %d", flag, size, math.MaxInt32)
	}
	return int32(size), nil
}
