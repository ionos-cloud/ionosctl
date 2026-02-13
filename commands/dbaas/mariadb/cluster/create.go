package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/spf13/cobra"
)

func Create() *core.Command {
	baseReqFlags := []string{
		constants.FlagName, constants.FlagVersion, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
		constants.ArgUser, constants.ArgPassword,
	}
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create DBaaS MariaDB clusters",
		Example:   fmt.Sprintf("i db mariadb cluster create %s", core.FlagsUsage(baseReqFlags...)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, baseReqFlags...)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := mariadb.CreateClusterProperties{}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				cluster.DisplayName = name
			}

			if c.Command.Command.Flags().Changed(constants.FlagInstances) {
				instances, err := c.Command.Command.Flags().GetInt32(constants.FlagInstances)
				if err != nil {
					return err
				}
				cluster.Instances = instances
			}

			if c.Command.Command.Flags().Changed(constants.FlagVersion) {
				version, err := c.Command.Command.Flags().GetString(constants.FlagVersion)
				if err != nil {
					return err
				}
				cluster.MariadbVersion = mariadb.MariadbVersion(version)
			}

			if c.Command.Command.Flags().Changed(constants.FlagCores) {
				cores, err := c.Command.Command.Flags().GetInt32(constants.FlagCores)
				if err != nil {
					return err
				}
				cluster.Cores = cores
			}

			if c.Command.Command.Flags().Changed(constants.FlagStorageSize) {
				storageSizeStr, err := c.Command.Command.Flags().GetString(constants.FlagStorageSize)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(storageSizeStr, convbytes.GB)
				cluster.StorageSize = int32(sizeInt64)
			}

			if c.Command.Command.Flags().Changed(constants.FlagRam) {
				ramStr, err := c.Command.Command.Flags().GetString(constants.FlagRam)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(ramStr, convbytes.GB)
				cluster.Ram = int32(sizeInt64)
			}

			cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceDay) {
				day, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceDay)
				if err != nil {
					return err
				}
				cluster.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(day)
			}

			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceTime) {
				maintenanceTime, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceTime)
				if err != nil {
					return err
				}
				cluster.MaintenanceWindow.Time = maintenanceTime
			}

			cluster.Connections = make([]mariadb.Connection, 1)
			if c.Command.Command.Flags().Changed(constants.FlagCidr) {
				cidr, err := c.Command.Command.Flags().GetString(constants.FlagCidr)
				if err != nil {
					return err
				}
				cluster.Connections[0].Cidr = cidr
			}

			if c.Command.Command.Flags().Changed(constants.FlagDatacenterId) {
				datacenterId, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
				if err != nil {
					return err
				}
				cluster.Connections[0].DatacenterId = datacenterId
			}

			if c.Command.Command.Flags().Changed(constants.FlagLanId) {
				lanId, err := c.Command.Command.Flags().GetString(constants.FlagLanId)
				if err != nil {
					return err
				}
				cluster.Connections[0].LanId = lanId
			}

			cluster.Credentials = mariadb.DBUser{}
			if c.Command.Command.Flags().Changed(constants.ArgUser) {
				username, err := c.Command.Command.Flags().GetString(constants.ArgUser)
				if err != nil {
					return err
				}
				cluster.Credentials.Username = username
			}

			if c.Command.Command.Flags().Changed(constants.ArgPassword) {
				password, err := c.Command.Command.Flags().GetString(constants.ArgPassword)
				if err != nil {
					return err
				}
				cluster.Credentials.Password = password
			}

			createdCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersPost(context.Background()).CreateClusterRequest(
				mariadb.CreateClusterRequest{Properties: &cluster},
			).Execute()
			if err != nil {
				return fmt.Errorf("failed creating cluster: %w", err)
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "10.6", "The MariaDB version of your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10.6", "10.11"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances of the cluster (one primary and n-1 secondaries)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagCores, "", 1, "Core count")
	cmd.AddStringFlag(constants.FlagRam, "", "4GB", "RAM size. e.g.: --ram 4GB. Minimum of 4GB. The maximum RAM size is determined by your contract limit")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "12GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	// Maintenance
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7) // Random hour 10-16
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", fmt.Sprintf("%02d:00:00", hour),
		"Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59. "+
			"Defaults to a random day during Mon-Fri, during the hours 10:00-16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))],
		"Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. "+
			"Defaults to a random day during Mon-Fri, during the hours 10:00-16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})
	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		datacenterId, _ := c.Flags().GetString(constants.FlagDatacenterId)
		return cloudapiv6completer.LansIds(datacenterId),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagCidr, "", "", "The IP and subnet for your cluster. All IPs must be in a /24 network", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCidr, completer.GetCidrCompletionFunc(cmd))
	// credentials / DBUser
	cmd.AddStringFlag(constants.ArgUser, "", "", "The initial username", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.ArgPassword, "", "", "The password", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
