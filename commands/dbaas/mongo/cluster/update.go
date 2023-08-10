package cluster

import (
	"context"
	"os"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

func ClusterUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a Mongo Cluster by ID",
		Example:   "ionosctl dbaas mongo cluster update --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr)
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			c.Printer.Verbose("Getting Cluster by id: %s", clusterId)

			updateProperties := ionoscloud.PatchClusterProperties{}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
				updateProperties.SetDisplayName(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagTemplateId)) {
				updateProperties.SetTemplateID(viper.GetString(core.GetFlagName(c.NS, constants.FlagTemplateId)))
			}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
				updateProperties.SetInstances(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)))
			}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) &&
				viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
				maintenanceWindow := ionoscloud.MaintenanceWindow{}
				maintenanceWindow.SetTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)))
				maintenanceWindow.SetDayOfTheWeek(ionoscloud.DayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))))
				updateProperties.SetMaintenanceWindow(maintenanceWindow)
			}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) &&
				viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) &&
				viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
				conn := ionoscloud.Connection{}
				conn.SetDatacenterId(viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId)))
				conn.SetLanId(viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId)))
				conn.SetCidrList(viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagCidr)))
				updateProperties.SetConnections([]ionoscloud.Connection{conn})
			}

			cluster, _, err := client.Must().MongoClient.ClustersApi.ClustersPatch(context.Background(), clusterId).
				PatchClusterRequest(
					ionoscloud.PatchClusterRequest{Properties: &updateProperties},
				).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, &[]ionoscloud.ClusterResponse{cluster}))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "6.0", "The MongoDB version of your cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "The physical location where the cluster will be created. (defaults to the location of the connected datacenter)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances in the replicaset cluster (one primary and n-1 secondaries). Setting this flag infers a replicaset type. Limited to at least 3 for business edition. (required for non-playground replicaset clusters)", core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagShards, "", 1, "The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters)", core.RequiredFlagOption())

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the Maintenance. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day for Maintenance. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Enterprise-specific
	cmd.AddIntFlag(constants.FlagCores, "", 0, "The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagRam, "", "", "Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (required and only settable for enterprise edition)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024MB", "2GB", "4GB", "8GB", "12GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagStorageType, "", "", []string{"HDD", "SSD", "\"SSD Premium\""},
		"Custom Storage Type. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagStorageSize, "", "", "Custom Storage: Greater performance for values greater than 100 GB. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"50GB", "100GB", "200GB", "1TB", "10TB", "100TB", "1000TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.FlagCidr, "", nil, "The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24", core.RequiredFlagOption())

	cmd.AddStringFlag(flagBackupLocation, "", "", "The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster")

	cmd.Command.SilenceUsage = true

	return cmd
}
