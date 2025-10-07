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
	"github.com/spf13/viper"
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
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.GetString(fn) != "" {
				cluster.Instances = viper.GetInt32(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				cluster.MariadbVersion = mariadb.MariadbVersion(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.GetString(fn) != "" {
				cluster.Cores = viper.GetInt32(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				cluster.StorageSize = int32(sizeInt64)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				cluster.Ram = int32(sizeInt64)
			}

			cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.Time = viper.GetString(fn)
			}

			cluster.Connections = make([]mariadb.Connection, 1)
			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				cluster.Connections[0].Cidr = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				cluster.Connections[0].DatacenterId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				cluster.Connections[0].LanId = viper.GetString(fn)
			}

			cluster.Credentials = mariadb.DBUser{}
			if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) {
				cluster.Credentials.Username = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
				cluster.Credentials.Password = viper.GetString(fn)
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
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
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
