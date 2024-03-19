package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	sdkgo "github.com/avirtopeanu-ionos/alpha-sdk-go-dbaas-mariadb"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*
   - instances can only be increased (3, 5, 7),
   - storageSize can only be increased,
   - ram and cores can be both increased and decreased.
*/

func Update() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a MariaDB Cluster",
		LongDesc: `This will update the MariaDB Cluster with the provided parameters.
The following parameters can be updated:
- instances can only be increased (3, 5, 7),
- storageSize can only be increased,
- ram and cores can be both increased and decreased.`,
		Example: "ionosctl dbaas mariadb cluster update --cluster-id CLUSTER_ID --cores 4 --ram 16",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {

			cluster := sdkgo.CreateClusterProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.GetString(fn) != "" {
				cluster.Instances = pointer.From(viper.GetInt32(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.GetString(fn) != "" {
				cluster.Cores = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.MB)
				cluster.StorageSize = pointer.From(int32(sizeInt64))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.MB)
				cluster.Ram = pointer.From(int32(sizeInt64))
			}

			cluster.MaintenanceWindow = &sdkgo.MaintenanceWindow{}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.DayOfTheWeek = (*sdkgo.DayOfTheWeek)(pointer.From(
					viper.GetString(fn)))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.Time = pointer.From(
					viper.GetString(fn))
			}

			createdCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersPost(context.Background()).CreateClusterRequest(
				sdkgo.CreateClusterRequest{Properties: &cluster},
			).Execute()
			if err != nil {
				return fmt.Errorf("failed creating cluster: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasMariadbCluster,
				createdCluster, tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ClustersProperty(func(c sdkgo.ClusterResponse) string {
			if c.Id == nil {
				return ""
			}
			return *c.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster")

	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances of the cluster (one primary and n-1 secondaries)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagCores, "", 1, "Core count")
	cmd.AddStringFlag(constants.FlagRam, "", "2GB", "Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024MB", "2GB", "4GB", "8GB", "12GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
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
		return append(workingDaysOfWeek, "Satuday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
