package cluster

import (
	"context"
	"os"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

func ClusterUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "dbaas-mongo",
			Resource:  "cluster",
			Verb:      "update",
			Aliases:   []string{"u"},
			ShortDesc: "Update a Mongo Cluster by ID",
			Example:   "ionosctl dbaas mongo cluster update --cluster-id <cluster-id>",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
				c.Command.Command.MarkFlagsRequiredTogether(
					constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
				)
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
					maintenanceWindow.SetDayOfTheWeek(
						ionoscloud.DayOfTheWeek(
							viper.GetString(
								core.GetFlagName(
									c.NS, constants.FlagMaintenanceDay,
								),
							),
						),
					)
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

				cluster, _, err := c.DbaasMongoServices.Clusters().Update(
					clusterId, ionoscloud.PatchClusterRequest{Properties: &updateProperties},
				)
				if err != nil {
					return err
				}
				return c.Printer.Print(getClusterPrint(c, &[]ionoscloud.ClusterResponse{cluster}))
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "When using text output, don't print headers")
	cmd.AddStringFlag(
		constants.FlagTemplateId, "", "",
		"The unique ID of the template, which specifies the number of cores, storage size, and memory. You cannot downgrade to a smaller template or minor edition",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagTemplateId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.MongoTemplateIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddInt32Flag(
		constants.FlagInstances, "", 0, "The total number of instances in the cluster (one primary and n-1 secondaries)",
	)

	// Maintenance
	cmd.AddStringFlag(
		constants.FlagMaintenanceTime, "", "",
		"Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59",
		core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(
		constants.FlagMaintenanceDay, "", "",
		"Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur",
		core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagMaintenanceDay,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	) // TODO: Completions should be a flag option func

	// Connections
	cmd.AddStringFlag(
		constants.FlagDatacenterId, "", "",
		"The datacenter to which your cluster will be connected. Must be in the same location as the cluster",
		core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagDatacenterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(
		constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLanId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cloudapiv6completer.LansIds(
				os.Stderr, cmd.Flag(constants.FlagDatacenterId).Value.String(),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringSliceFlag(
		constants.FlagCidr, "", nil,
		"The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.114.0/24",
		core.RequiredFlagOption(),
	)

	cmd.Command.SilenceUsage = true

	return cmd
}
