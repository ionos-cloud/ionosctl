package cluster

import (
	"context"
	"fmt"
	"math"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/spf13/cobra"

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
		Example:   "ionosctl dbaas mariadb cluster update" + core.FlagsUsage(constants.ClusterId, constants.FlagVersion),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := mariadb.PatchClusterProperties{}

			if c.Command.Command.Flags().Changed(constants.FlagVersion) {
				version, err := c.Command.Command.Flags().GetString(constants.FlagVersion)
				if err != nil {
					return err
				}
				cluster.MariadbVersion = pointer.From(mariadb.MariadbVersion(version))
			}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				cluster.DisplayName = pointer.From(name)
			}

			if c.Command.Command.Flags().Changed(constants.FlagInstances) {
				instances, err := c.Command.Command.Flags().GetInt32(constants.FlagInstances)
				if err != nil {
					return err
				}
				cluster.Instances = pointer.From(instances)
			}

			if c.Command.Command.Flags().Changed(constants.FlagCores) {
				cores, err := c.Command.Command.Flags().GetInt32(constants.FlagCores)
				if err != nil {
					return err
				}
				cluster.Cores = pointer.From(cores)
			}

			if c.Command.Command.Flags().Changed(constants.FlagStorageSize) {
				storageSizeStr, err := c.Command.Command.Flags().GetString(constants.FlagStorageSize)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(storageSizeStr, convbytes.GB)
				if sizeInt64 < 0 || sizeInt64 > math.MaxInt32 {
					return fmt.Errorf("storage size %d is out of allowed int32 range [0-%d]", sizeInt64, math.MaxInt32)
				}
				cluster.StorageSize = pointer.From(int32(sizeInt64))
			}

			if c.Command.Command.Flags().Changed(constants.FlagRam) {
				ramStr, err := c.Command.Command.Flags().GetString(constants.FlagRam)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(ramStr, convbytes.GB)
				if sizeInt64 < 0 || sizeInt64 > math.MaxInt32 {
					return fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", sizeInt64, math.MaxInt32)
				}
				cluster.Ram = pointer.From(int32(sizeInt64))
			}

			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceDay) {
				day, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceDay)
				if err != nil {
					return err
				}
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(day)
			}

			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceTime) {
				maintenanceTime, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceTime)
				if err != nil {
					return err
				}
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.Time = maintenanceTime
			}

			clusterId, err := c.Command.Command.Flags().GetString(constants.FlagClusterId)
			if err != nil {
				return err
			}
			createdCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersPatch(context.Background(), clusterId).
				PatchClusterRequest(mariadb.PatchClusterRequest{Properties: &cluster}).Execute()
			if err != nil {
				return fmt.Errorf("failed updating cluster: %w", err)
			}

			converted, err := resource2table.ConvertDbaasMariaDBClusterToTable(createdCluster)
			if err != nil {
				return fmt.Errorf("failed converting cluster to table: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(
				createdCluster,
				converted,
				cols,
			)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster",
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster")
	cmd.AddInt32Flag(constants.FlagInstances, "", 0, "The total number of instances of the cluster (one primary and n-1 secondaries). Instances can only be increased (3,5,7)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagVersion, "", "", "The MariaDB version of your cluster. Downgrades are not supported (version can only be increased) ")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10.6", "10.11"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddInt32Flag(constants.FlagCores, "", 0, "Core count. Can be increased or decreased.")
	cmd.AddStringFlag(constants.FlagRam, "", "", "RAM size. e.g.: --ram 4GB. Can be increased or decreased.")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "12GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", "", "The size of the Storage in GB. Can only be increased")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. e.g.: 16:30:59. To change maintenance provide both --maintenance-day and --maintenance-time")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. e.g.: Monday. To change maintenance provide both --maintenance-day and --maintenance-time")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// They do nothing... but we can't outright remove them in case some user already uses them in their scripts
	// would cause ('unknown flag: -w')
	_ = cmd.Command.Flags().MarkHidden(constants.ArgWaitForRequest)
	_ = cmd.Command.Flags().MarkHidden(constants.ArgTimeout)

	cmd.Command.SilenceUsage = true

	return cmd
}
