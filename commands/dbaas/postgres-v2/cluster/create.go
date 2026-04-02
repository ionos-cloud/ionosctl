package cluster

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCreateCmd() *core.Command {
	ctx := context.TODO()

	// Generate random maintenance window defaults (Mon-Fri, 10:00-16:00)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7)
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}
	defaultMaintenanceDay := workingDaysOfWeek[r.Intn(len(workingDaysOfWeek))]
	defaultMaintenanceTime := fmt.Sprintf("%02d:00:00", hour)

	create := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL Cluster",
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* PostgreSQL Version
* Credentials for the database user: Username, Password, and Database name`,
		Example:    "ionosctl dbaas postgres-v2 cluster create --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --db-username <username> --db-password <password> --database <database> --version <version>",
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "", "The PostgreSQL version of your Cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 1, "The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5")
	create.AddIntFlag(constants.FlagCores, "", 2, "The number of CPU cores per instance. Minimum: 1")
	create.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance in GB. e.g. --ram 4096, --ram 4096MB, --ram 4GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "de",
		"The S3 location where the backups will be stored. Defaults to 'de'")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		locations, _, err := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
		if err != nil {
			return []string{"de"}, cobra.ShellCompDirectiveNoFileComp
		}
		return functional.Map(locations.Items, func(l psqlv2.BackupLocationRead) string {
			if l.Properties.Location != nil {
				return *l.Properties.Location
			}
			return l.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagSyncMode, constants.FlagSyncModeShort, "ASYNCHRONOUS", "Replication mode: ASYNCHRONOUS, STRICTLY_SYNCHRONOUS")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "20GB", "The amount of storage per instance in GB. e.g.: --storage-size 20480 or --storage-size 20480MB or --storage-size 20GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, constants.FlagLanIdShortPsql, "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCidr, constants.FlagCidrShortPsql, "", "The IP and subnet for the cluster. Note the following unavailable IP range: 10.208.0.0/12. e.g.: 192.168.1.100/24", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagBackupId, constants.FlagBackupIdShortPsql, "", "The unique ID of the backup you want to restore from when creating this cluster")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		backups, err := backup.Backups()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		const timeFmt = "2006-01-02 15:04"
		return functional.Map(backups.Items, func(b psqlv2.BackupRead) string {
			latest := "now"
			if b.Properties.LatestRecoveryTargetTime != nil {
				latest = b.Properties.LatestRecoveryTargetTime.Time.Format(timeFmt)
			}
			earliest := "n/a"
			if b.Properties.EarliestRecoveryTargetTime != nil {
				earliest = b.Properties.EarliestRecoveryTargetTime.Time.Format(timeFmt)
			}
			return fmt.Sprintf("%s\tfor cluster '%s': earliest: '%s', latest: '%s'",
				b.Id, *b.Properties.ClusterId, earliest, latest)
		}), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")

	create.AddStringFlag(constants.FlagDbUsername, constants.FlagDbUsernameShortPsql, "", "Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby)", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDbPassword, constants.FlagDbPasswordShortPsql, "", "Password for the initial postgres user", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDatabase, "", "", "The name of the initial database to be created", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")
	create.AddStringFlag(constants.FlagConnectionPooler, "", "DISABLED", "Connection pooling mode: DISABLED, TRANSACTION, SESSION")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagConnectionPooler, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"DISABLED", "TRANSACTION", "SESSION"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Enable collection and reporting of logs for this cluster")
	create.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Enable collection and reporting of metrics for this cluster")

	create.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, defaultMaintenanceTime,
		"Time for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59. "+
			"Defaults to a random time during 10:00-16:00")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, defaultMaintenanceDay,
		"Day of the week for the MaintenanceWindow. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. "+
			"Defaults to a random day during Mon-Fri")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	create.AddStringSliceFlag(constants.ArgCols, "", table.DefaultCols(clusterCols), table.ColsMessage(clusterCols))
	_ = create.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(clusterCols), cobra.ShellCompDirectiveNoFileComp
	})

	return create
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr, constants.FlagDbUsername, constants.FlagDbPassword, constants.FlagDatabase, constants.FlagVersion)
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

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating Cluster...")

	cluster, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersPost(context.Background()).ClusterCreate(input).Execute()
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := cluster.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, *id); err != nil {
				return err
			}

			if cluster, _, err = client.Must().PostgresClientV2.ClustersApi.
				ClustersFindById(context.Background(), *id).Execute(); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new Cluster Id")
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}

func getCreateClusterRequest(c *core.CommandConfig) (psqlv2.ClusterCreate, error) {
	inputCluster := psqlv2.ClusterCreate{}
	input := psqlv2.ClusterCreateProperties{}

	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	c.Verbose("PostgresVersion: %v", pgsqlVersion)
	input.SetVersion(pgsqlVersion)

	instanceConfig := psqlv2.InstanceConfiguration{}
	replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
	c.Verbose("Instances: %v", replicas)
	instanceConfig.Count = replicas

	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	c.Verbose("Cores: %v", cpuCoreCount)
	instanceConfig.Cores = cpuCoreCount

	// Convert Ram (SDK expects GB)
	size, ok := convbytes.StrToUnitOk(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), convbytes.GB)
	if !ok {
		return inputCluster, fmt.Errorf("invalid value for Ram: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)))
	}
	if size < 0 || size > math.MaxInt32 {
		return inputCluster, fmt.Errorf("Ram value %vGB exceeds valid range", size)
	}
	instanceConfig.Ram = int32(size)
	c.Verbose("Ram: %vGB", int32(size))

	// Convert StorageSize (SDK expects GB)
	storageSize, ok := convbytes.StrToUnitOk(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), convbytes.GB)
	if !ok {
		return inputCluster, fmt.Errorf("invalid value for StorageSize: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)))
	}
	if storageSize < 0 || storageSize > math.MaxInt32 {
		return inputCluster, fmt.Errorf("StorageSize value %vGB exceeds valid range", storageSize)
	}
	instanceConfig.SetStorageSize(int32(storageSize))
	c.Verbose("StorageSize: %vGB", int32(storageSize))

	input.SetInstances(instanceConfig)

	// BackupLocation is required - default is "de"
	backupLoc := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation))
	if backupLoc == "" {
		backupLoc = "de"
	}
	c.Verbose("BackupLocation: %v", backupLoc)
	input.SetBackupLocation(backupLoc)

	displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	c.Verbose("DisplayName: %v", displayName)
	input.SetName(displayName)

	dbuser := psqlv2.PostgresUser{}
	username := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbUsername))
	c.Verbose("DBUser - Username: %v", username)
	dbuser.SetUsername(username)

	password := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbPassword))
	dbuser.SetPassword(password)

	database := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))
	c.Verbose("DBUser - Database: %v", database)
	dbuser.SetDatabase(database)

	input.SetCredentials(dbuser)

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
		desc := viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))
		c.Verbose("Description: %v", desc)
		input.SetDescription(desc)
	}

	connectionPooler := viper.GetString(core.GetFlagName(c.NS, constants.FlagConnectionPooler))
	c.Verbose("ConnectionPooler: %v", connectionPooler)
	input.SetConnectionPooler(connectionPooler)

	logsEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagLogsEnabled))
	c.Verbose("LogsEnabled: %v", logsEnabled)
	input.SetLogsEnabled(logsEnabled)

	metricsEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagMetricsEnabled))
	c.Verbose("MetricsEnabled: %v", metricsEnabled)
	input.SetMetricsEnabled(metricsEnabled)

	vdcConnection := psqlv2.PostgresClusterConnection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
	c.Verbose("Connection - DatacenterId: %v", vdcId)
	vdcConnection.SetDatacenterId(vdcId)

	lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
	c.Verbose("Connection - LanId: %v", lanId)
	vdcConnection.SetLanId(lanId)

	ip := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
	c.Verbose("Connection - Cidr: %v", ip)
	vdcConnection.SetPrimaryInstanceAddress(ip)

	input.SetConnection(vdcConnection)

	// MaintenanceWindow is required - always set from flags (which have defaults)
	maintenanceWindow := psqlv2.MaintenanceWindow{}
	maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
	c.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
	maintenanceWindow.SetTime(maintenanceTime)

	maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
	c.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
	maintenanceWindow.SetDayOfTheWeek(psqlv2.DayOfTheWeek(maintenanceDay))

	input.SetMaintenanceWindow(maintenanceWindow)

	// ReplicationMode is required - set from sync-mode flag
	syncMode := viper.GetString(core.GetFlagName(c.NS, constants.FlagSyncMode))
	c.Verbose("ReplicationMode: %v", syncMode)
	input.SetReplicationMode(psqlv2.PostgresClusterReplicationMode(syncMode))

	// Restore from backup (optional)
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
		backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
		restoreFromBackup := psqlv2.NewPostgresClusterFromBackup(backupId)

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
			recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return inputCluster, fmt.Errorf("invalid recovery-time format (expected RFC3339, e.g. 2024-01-15T10:00:00Z): %w", err)
			}

			c.Verbose("From Backup - RecoveryTargetTime: %v", recoveryTargetTime)
			targetTime := psqlv2.IonosTime{Time: recoveryTargetTime}
			restoreFromBackup.RecoveryTargetDatetime = &targetTime
		}

		if backupId != "" {
			c.Verbose("From Backup - BackupId: %v", backupId)
		}

		input.RestoreFromBackup = restoreFromBackup
	}

	inputCluster.SetProperties(input)
	return inputCluster, nil
}
