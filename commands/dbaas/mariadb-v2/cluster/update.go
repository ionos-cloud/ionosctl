package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterUpdateCmd() *core.Command {
	ctx := context.TODO()

	update := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb-v2",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a MariaDB Cluster",
		LongDesc: `Use this command to update attributes of a MariaDB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

Note the API's sizing constraints: instances and storage size can only be increased (never decreased), the version can only be upgraded (no downgrade), while cores and RAM can be both increased and decreased.

Required values to run command:

* Cluster Id`,
		Example: "ionosctl dbaas mariadb-v2 cluster update --cluster-id <cluster-id> --cores 4 --ram 8GB",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays)) {
				if err := validateBackupRetentionDays(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays))); err != nil {
					return err
				}
			}
			return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId)
		},
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	update.AddStringFlag(constants.FlagVersion, "", "", "The MariaDB version of your cluster. Downgrades are not supported",
		core.WithCompletion(completer.Versions, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	update.AddInt32Flag(constants.FlagInstances, "", 0, "The total number of instances in the cluster (one primary and n-1 secondary). Can only be increased")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddInt32Flag(constants.FlagCores, "", 0, "The number of CPU cores per instance. Can be increased or decreased")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory per instance. e.g. --ram 4, --ram 4GB. Can be increased or decreased")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "The size of the storage per instance. e.g. --storage-size 10GB. Can only be increased")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The friendly name of your cluster")
	update.AddInt32Flag(constants.FlagBackupRetentionDays, "", 0, "Configures how many days cluster backups are retained. Minimum: 1, Maximum: 365")
	update.AddBoolFlag(constants.FlagLogsEnabled, "", false, "Enable collection and reporting of logs for this cluster")
	update.AddBoolFlag(constants.FlagMetricsEnabled, "", false, "Enable collection and reporting of metrics for this cluster")
	update.AddStringFlag(constants.FlagMaintenanceTime, "", "",
		"Time for the MaintenanceWindow. e.g.: 16:30:59. Must be specified together with --maintenance-day")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddSetFlag(constants.FlagMaintenanceDay, "", "",
		[]string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		"Day of the week for the MaintenanceWindow. Must be specified together with --maintenance-time")

	// Credentials are only re-sent when --password is given (see updateClusterProperties).
	update.AddStringFlag(constants.ArgUser, "", "", "New username for the database user. Only applied together with --password")
	update.AddStringFlag(constants.ArgPassword, "", "", "New password for the database user. The API does not return credentials on GET, so both --user and --database must be supplied alongside it")
	update.AddStringFlag(constants.FlagDatabase, "", "", "Database for the credentials. Only applied together with --password")

	return update
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	c.Verbose(constants.ClusterId, clusterId)

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().MariaClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	// Warn before a version change: backups taken on the current version cannot be
	// restored onto the upgraded cluster (the API requires an exact
	// backup-vs-cluster version match), so an upgrade effectively strands
	// pre-upgrade backups for in-place restore.
	if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
		newVersion := viper.GetString(fn)
		currentVersion := clusterRead.Properties.Version
		if newVersion != "" && newVersion != currentVersion {
			prompt := fmt.Sprintf(
				"changing version %s -> %s: backups taken on %s can no longer be restored in place onto this cluster "+
					"(the API requires the backup and cluster versions to match). To recover pre-upgrade data later, create a NEW "+
					"cluster from the old backup at version %s (`cluster create --backup-id <id> --version %s`) and then upgrade it. Continue",
				currentVersion, newVersion, currentVersion, currentVersion, currentVersion)
			if !confirm.FAsk(c.Command.Command.InOrStdin(), prompt, viper.GetBool(constants.ArgForce)) {
				return fmt.Errorf(confirm.UserDenied)
			}
		}
	}

	newCluster, err := updateClusterProperties(c, clusterRead.Properties)
	if err != nil {
		return err
	}

	// Update (Ensure) Cluster
	c.Verbose("Updating Cluster...")
	clusterEnsure := mariadb.NewClusterEnsure(clusterId, newCluster)

	item, _, err := client.Must().MariaClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(clusterCols, item, cols))
}

func updateClusterProperties(c *core.CommandConfig, input mariadb.Cluster) (mariadb.Cluster, error) {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
		instances := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
		c.Verbose("Instances: %v", instances)
		input.Instances.Count = instances
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		c.Verbose("Cores: %v", cores)
		input.Instances.Cores = cores
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		ram, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), "ram")
		if err != nil {
			return input, err
		}
		c.Verbose("Ram: %vGB", ram)
		input.Instances.Ram = ram
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStorageSize)) {
		storage, err := toInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), "storage-size")
		if err != nil {
			return input, err
		}
		c.Verbose("StorageSize: %vGB", storage)
		input.Instances.StorageSize = storage
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		version := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		c.Verbose("Version: %v", version)
		input.Version = version
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
		c.Verbose("Name: %v", name)
		input.Name = name
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays)) {
		retention := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagBackupRetentionDays))
		c.Verbose("BackupRetentionDays: %v", retention)
		input.Backup.RetentionDays = retention
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
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
		input.MaintenanceWindow.Time = viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		input.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)))
	}

	// The API does not return credentials on GET, so they are only re-sent when the
	// user supplies a new --password; otherwise the existing credentials are kept.
	if err := applyCredentialsFromFlags(c, &input); err != nil {
		return input, err
	}

	return input, nil
}
