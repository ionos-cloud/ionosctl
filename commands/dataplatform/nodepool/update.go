package nodepool

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/dataplatform/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolUpdateCmd() *core.Command {
	var (
		updateProperties = dataplatform.PatchNodePoolProperties{}
	)

	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "update", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"u"},
		ShortDesc: "Update Dataplatform Nodepools",
		LongDesc:  "Node pools are the resources that powers the DataPlatformCluster.\n\nThe following requests allows to alter the existing resources of the cluster",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			err = c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagNodepoolId)
			if err != nil {
				return err
			}
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating Nodepool..."))

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			npId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

			fd := core.GetFlagName(c.NS, constants.FlagMaintenanceDay)
			ft := core.GetFlagName(c.NS, constants.FlagMaintenanceTime)
			if viper.IsSet(fd) && viper.IsSet(ft) {
				maintenanceWindow := dataplatform.MaintenanceWindow{}
				maintenanceWindow.SetDayOfTheWeek(viper.GetString(fd))
				maintenanceWindow.SetTime(viper.GetString(ft))
				updateProperties.SetMaintenanceWindow(maintenanceWindow)
			}
			if f := core.GetFlagName(c.NS, constants.FlagNodeCount); viper.IsSet(f) {
				updateProperties.SetNodeCount(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagAnnotations); viper.IsSet(f) {
				updateProperties.SetAnnotations(viper.GetStringMap(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagLabels); viper.IsSet(f) {
				updateProperties.SetLabels(viper.GetStringMap(f))
			}

			input := dataplatform.PatchNodePoolRequest{}
			input.SetProperties(updateProperties)

			cr, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsPatch(context.Background(), clusterId, npId).PatchNodePoolRequest(input).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			crConverted, err := resource2table.ConvertDataplatformNodePoolToTable(cr)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(cr, crConverted, tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagClusterId, "", "", "The UUID of the cluster the nodepool belongs to")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(constants.FlagNodepoolId, constants.FlagIdShort, "", "The UUID of the cluster the nodepool belongs to")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagNodepoolId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformNodepoolsIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagNodeCount, "n", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request [seconds]")

	cmd.Command.SilenceUsage = true

	return cmd
}
