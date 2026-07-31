package cluster

import (
	"context"
	"fmt"
	"math"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
)

// Update creates the `dbaas-mariadb cluster update` command which allows updating
// a MariaDB cluster. Notes for users are included in the help text: instances can
// only be increased (3, 5, 7), mariadbVersion can only be increased (no downgrade),
// storageSize can only be increased, ram and cores can be both increased and decreased.
func Update() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a MariaDB Cluster",
		LongDesc: `Update mutable attributes of an existing MariaDB cluster. Only the flags you pass are changed; everything else is left untouched. The cluster must be in state AVAILABLE for the update to be accepted.

Some changes are one-directional and cannot be reverted:
  - --version can only be upgraded (10.6 -> 10.11), never downgraded.
  - --instances can only be increased, and must stay odd (1 -> 3 -> 5).
  - --storage-size can only be increased, never shrunk.
  - --cores and --ram can be scaled up or down.

--maintenance-day and --maintenance-time must be supplied together (a maintenance window has both a day and a start time). The connection (datacenter/LAN/CIDR) and initial credentials are set at creation and cannot be changed here.`,
		Example: `# Upgrade the MariaDB version
ionosctl dbaas mariadb cluster update --cluster-id <cluster-id> --version 10.11

# Scale compute and move the maintenance window (day and time must be given together)
ionosctl dbaas mariadb cluster update --cluster-id <cluster-id> --cores 8 --ram 32GB --storage-size 200GB --maintenance-day Saturday --maintenance-time 02:00:00`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := mariadb.PatchClusterProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				cluster.MariadbVersion = pointer.From(mariadb.MariadbVersion(viper.GetString(fn)))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.IsSet(fn) {
				cluster.Instances = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(fn) {
				cluster.Cores = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.IsSet(fn) {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				if sizeInt64 < 0 || sizeInt64 > math.MaxInt32 {
					return fmt.Errorf("storage size %d is out of allowed int32 range [0-%d]", sizeInt64, math.MaxInt32)
				}
				cluster.StorageSize = pointer.From(int32(sizeInt64))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.IsSet(fn) {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				if sizeInt64 < 0 || sizeInt64 > math.MaxInt32 {
					return fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", sizeInt64, math.MaxInt32)
				}
				cluster.Ram = pointer.From(int32(sizeInt64))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.IsSet(fn) {
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.IsSet(fn) {
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.Time = viper.GetString(fn)
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			createdCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersPatch(context.Background(), clusterId).
				PatchClusterRequest(mariadb.PatchClusterRequest{Properties: &cluster}).Execute()
			if err != nil {
				return fmt.Errorf("failed updating cluster: %w", err)
			}

			return c.Printer(allCols).Print(createdCluster)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster to update. Must be in state AVAILABLE",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return ClustersProperty(func(c mariadb.ClusterResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id
				})
			}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "New human-friendly display name for the cluster (max 63 characters)")
	cmd.AddInt32Flag(constants.FlagInstances, "", 0, "New instance count (primary + secondaries). Can only be increased and must stay odd: 1, 3 or 5. Adding instances converts a standalone cluster into a high-availability replica set")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagVersion, "", "", "Upgrade the MariaDB version (one of: 10.6, 10.11). Upgrade only; downgrades are rejected by the API")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10.6", "10.11"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddInt32Flag(constants.FlagCores, "", 0, "New CPU core count per instance (minimum 1). Can be scaled up or down")
	cmd.AddStringFlag(constants.FlagRam, "", "", "New memory per instance, e.g. --ram 8GB. Minimum 4GB, whole GB only. Can be scaled up or down")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "12GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", "", "New storage per instance, e.g. --storage-size 200GB. Can only be increased (never shrunk), up to 2000GB (2TB)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "New start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window. Must be supplied together with --maintenance-day")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "New day of the week for the weekly maintenance window, e.g. Monday. Must be supplied together with --maintenance-time")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	return cmd
}
