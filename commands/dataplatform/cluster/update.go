package cluster

import (
	"context"
	"fmt"
	"os"

	"github.com/cilium/fake"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	sdkdataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
)

func ClusterUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a Dataplatform Cluster by ID",
		LongDesc:  "Modifies the specified DataPlatformCluster by its distinct cluster ID. The fields in the request body are applied to the cluster. Note that the application to the cluster itself is performed asynchronously. You can check the sync state by querying the cluster with the GET method",
		Example:   "ionosctl dataplatform cluster update --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr)
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Cluster by id: %s", clusterId))

			updateProperties := sdkdataplatform.PatchClusterProperties{}
			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
				updateProperties.SetName(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
				updateProperties.SetDataPlatformVersion(viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion)))
			}

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) &&
				viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
				maintenanceWindow := sdkdataplatform.MaintenanceWindow{}
				maintenanceWindow.SetTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)))
				maintenanceWindow.SetDayOfTheWeek(viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)))
				updateProperties.SetMaintenanceWindow(maintenanceWindow)
			}

			cluster, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersPatch(c.Context, clusterId).PatchClusterRequest(sdkdataplatform.PatchClusterRequest{Properties: &updateProperties}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, cluster,
				printer.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Stdout, out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fake.Names(10), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagVersion, "", "", "The version of the cluster")

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, cmd.Flag(constants.FlagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.FlagCidr, "", nil, "The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.114.0/24", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true

	return cmd
}
