package nodepool

import (
	"context"
	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	sdkdataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	annotationsBare      map[string]string
	labelsBare           map[string]string
	availabilityZoneBare string
	createProperties     = sdkdataplatform.CreateNodePoolProperties{}
)

func NodepoolCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create Dataplatform Nodepools",
		LongDesc:  "Node pools are the resources that powers the DataPlatformCluster.\n\nThe following requests allows to alter the existing resources, add or remove new resources to the cluster.",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			err = c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagNodeCount)
			if err != nil {
				return err
			}
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Creating Nodepool...")
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			day := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			time := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))

			if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) &&
				viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
				maintenanceWindow := sdkdataplatform.MaintenanceWindow{}
				maintenanceWindow.SetDayOfTheWeek(day)
				maintenanceWindow.SetTime(time)
				createProperties.SetMaintenanceWindow(maintenanceWindow)
			}

			input := sdkdataplatform.CreateNodePoolRequest{}
			input.SetProperties(createProperties)

			client, err := config.GetClient()
			if err != nil {
				return err
			}
			cr, _, err := client.DataplatformClient.DataPlatformNodePoolApi.CreateClusterNodepool(context.Background(), clusterId).CreateNodePoolRequest(input).Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getNodepoolsPrint(c, &[]sdkdataplatform.NodePoolResponseData{cr}))
		},
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The UUID of the cluster the nodepool belongs to")

	// Linked to properties struct
	//_ = allocate.Zero(&createProperties)
	cmd.AddStringVarFlag(createProperties.Name, constants.FlagName, constants.FlagNameShort, "", "The name of your nodepool", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fake.Names(10), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32VarFlag(createProperties.NodeCount, constants.FlagNodeCount, "", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddInt32VarFlag(createProperties.CoresCount, constants.FlagCores, "", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddInt32VarFlag(createProperties.RamSize, constants.FlagRam, "", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddInt32VarFlag(createProperties.StorageSize, constants.FlagStorageSize, "", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddStringVarFlag((*string)(createProperties.StorageType), constants.FlagStorageType, "", "", "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddStringVarFlag((*string)(createProperties.AvailabilityZone), constants.FlagAvailabilityZone, "", "", "The number of nodes that make up the node pool", core.RequiredFlagOption())

	cmd.AddStringToStringVarFlag(&labelsBare, constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringVarFlag(&annotationsBare, constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")

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
