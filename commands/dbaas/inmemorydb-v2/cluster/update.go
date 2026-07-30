package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterUpdateCmd() *core.Command {
	ctx := context.TODO()

	update := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-inmemorydb-v2",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an In-Memory DB Cluster",
		LongDesc: `Use this command to update attributes of an In-Memory DB Cluster. This command uses a combination of GET and PUT to simulate a PATCH operation.

Required values to run command:

* Cluster Id`,
		Example: "ionosctl dbaas in-memory-db-v2 cluster update --cluster-id <cluster-id> --password <password> --cores 4 --ram 8GB",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId, constants.ArgPassword)
		},
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	update.AddStringFlag(constants.FlagVersion, "", "", "The In-Memory DB version of your cluster",
		core.WithCompletion(completer.Versions, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	update.AddIntFlag(constants.FlagReplicas, "", 0, "The total number of replicas in the cluster (one active and n-1 passive)")
	update.AddIntFlag(constants.FlagCores, "", 0, "The number of CPU cores per instance")
	update.AddStringFlag(constants.FlagRam, "", "", "The amount of memory per instance in GB. e.g. --ram 4, --ram 4GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB", "128GB", "256GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The friendly name of your cluster")
	update.AddStringFlag(constants.FlagDescription, "", "", "Human-readable description for the cluster")
	update.AddStringFlag(constants.ArgUser, "", "", "Username for the initial In-Memory DB user. Defaults to the cluster's current username")
	update.AddStringFlag(constants.ArgPassword, "", "", "Password for the In-Memory DB user. Required because the API does not return it on GET requests", core.RequiredFlagOption())
	update.AddBoolFlag(constants.ArgHashPassword, "", true, "Hash plaintext passwords before sending. The API only accepts hashed passwords")
	update.AddSetFlag(constants.FlagPersistenceMode, "", "", []string{"None", "AOF", "RDB", "RDB_AOF"}, "Specifies how and if data is persisted")
	update.AddSetFlag(constants.FlagEvictionPolicy, "", "",
		[]string{"noeviction", "allkeys-lru", "allkeys-lfu", "allkeys-random", "volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl"}, "The eviction policy for the cluster")
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
	return update
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	c.Verbose(constants.ClusterId, clusterId)

	// Fetch existing cluster
	c.Verbose("Getting Cluster...")
	clusterRead, _, err := client.Must().InMemoryDBClientV2.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}

	newCluster, err := updateClusterProperties(c, clusterRead.Properties)
	if err != nil {
		return err
	}

	// Update (Ensure) Cluster
	c.Verbose("Updating Cluster...")
	clusterEnsure := inmemorydb.NewClusterEnsure(clusterId, newCluster)

	item, _, err := client.Must().InMemoryDBClientV2.ClustersApi.
		ClustersPut(context.Background(), clusterId).
		ClusterEnsure(*clusterEnsure).
		Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	return c.Out(table.Sprint(clusterCols, item, cols))
}

func updateClusterProperties(c *core.CommandConfig, input inmemorydb.Cluster) (inmemorydb.Cluster, error) {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagReplicas)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagReplicas))
		c.Verbose("Replicas: %v", replicas)
		input.Instances.Count = replicas
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		c.Verbose("Cores: %v", cores)
		input.Instances.Cores = cores
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		ram, err := ramToInt32GB(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)))
		if err != nil {
			return input, err
		}
		c.Verbose("Ram: %vGB", ram)
		input.Instances.Ram = ram
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		version := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		c.Verbose("Version: %v", version)
		input.SetVersion(version)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
		c.Verbose("Name: %v", displayName)
		input.SetName(displayName)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDescription)) {
		desc := viper.GetString(core.GetFlagName(c.NS, constants.FlagDescription))
		c.Verbose("Description: %v", desc)
		input.SetDescription(desc)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagPersistenceMode)) {
		pm := inmemorydb.PersistenceMode(viper.GetString(core.GetFlagName(c.NS, constants.FlagPersistenceMode)))
		c.Verbose("PersistenceMode: %v", pm)
		input.PersistenceMode = &pm
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagEvictionPolicy)) {
		ep := inmemorydb.EvictionPolicy(viper.GetString(core.GetFlagName(c.NS, constants.FlagEvictionPolicy)))
		c.Verbose("EvictionPolicy: %v", ep)
		input.EvictionPolicy = ep
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
		input.MaintenanceWindow.SetTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		input.MaintenanceWindow.SetDayOfTheWeek(inmemorydb.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))))
	}

	// Credentials: the API does not return the password on GET, so the fetched
	// cluster carries a password with an empty algorithm that a PUT would reject.
	// Always rebuild credentials from the required --password (and keep the
	// existing username unless --user overrides it).
	credentials := inmemorydb.ClusterCredentials{}
	if input.Credentials != nil {
		credentials = *input.Credentials
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) {
		credentials.Username = viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	}
	if credentials.Username == "" {
		return input, fmt.Errorf("could not determine username from the existing cluster; pass --%s", constants.ArgUser)
	}
	credentials.Password = buildHashedPassword(
		viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword)),
		viper.GetBool(core.GetFlagName(c.NS, constants.ArgHashPassword)),
	)
	input.Credentials = &credentials
	c.Verbose("Credentials - Username: %v", credentials.Username)

	return input, nil
}
