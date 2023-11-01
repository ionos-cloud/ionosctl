package nodepool

import (
	"context"
	"fmt"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	sdkdataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolCreateCmd() *core.Command {
	var (
		// forward declaring required opts is fine, as we don't mind bad defaults
		name             string
		nodeCount        int32
		createProperties = sdkdataplatform.CreateNodePoolProperties{}
	)

	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create Dataplatform Nodepools",
		LongDesc:  "Node pools are the resources that powers the DataPlatformCluster.\n\nThe following requests allows to alter the existing resources, add or remove new resources to the cluster.",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			err = c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			if err != nil {
				return err
			}
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
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Creating Nodepool..."))
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			createProperties.Name = &name
			createProperties.NodeCount = &nodeCount

			fd := core.GetFlagName(c.NS, constants.FlagMaintenanceDay)
			ft := core.GetFlagName(c.NS, constants.FlagMaintenanceTime)
			if viper.IsSet(fd) && viper.IsSet(ft) {
				maintenanceWindow := sdkdataplatform.MaintenanceWindow{}
				maintenanceWindow.SetDayOfTheWeek(viper.GetString(fd))
				maintenanceWindow.SetTime(viper.GetString(ft))
				createProperties.SetMaintenanceWindow(maintenanceWindow)
			}
			if f := core.GetFlagName(c.NS, constants.FlagAnnotations); viper.IsSet(f) {
				createProperties.SetAnnotations(viper.GetStringMap(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagLabels); viper.IsSet(f) {
				createProperties.SetLabels(viper.GetStringMap(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(f) {
				createProperties.SetCoresCount(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagRam); viper.IsSet(f) {
				createProperties.SetRamSize(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(f) {
				createProperties.SetCpuFamily(viper.GetString(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.IsSet(f) {
				createProperties.SetStorageSize(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagStorageType); viper.IsSet(f) {
				createProperties.SetStorageType(sdkdataplatform.StorageType(viper.GetString(f)))
			}
			if f := core.GetFlagName(c.NS, constants.FlagAvailabilityZone); viper.IsSet(f) {
				createProperties.SetAvailabilityZone(sdkdataplatform.AvailabilityZone(f))
			}

			input := sdkdataplatform.CreateNodePoolRequest{}
			input.SetProperties(createProperties)

			cr, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsPost(context.Background(), clusterId).CreateNodePoolRequest(input).Execute()
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The UUID of the cluster the nodepool belongs to")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	// Linked to properties struct
	cmd.AddStringVarFlag(&name, constants.FlagName, constants.FlagNameShort, "", "The name of your nodepool", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fake.Names(10), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32VarFlag(&nodeCount, constants.FlagNodeCount, "", 0, "The number of nodes that make up the node pool", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagCpuFamily, "", "", "A valid CPU family name or AUTO if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource")
	cmd.AddInt32Flag(constants.FlagCores, "", 0, "The number of CPU cores per node")
	cmd.AddInt32Flag(constants.FlagRam, "", 0, "The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB")
	cmd.AddStringFlag(constants.FlagAvailabilityZone, "", "", "The availability zone of the virtual datacenter region where the node pool resources should be provisioned")
	cmd.AddStringFlag(constants.FlagStorageType, "", "", "The type of hardware for the volume")
	cmd.AddInt32Flag(constants.FlagStorageSize, "", 0, "The size of the volume in GB. The size must be greater than 10GB")

	cmd.AddStringToStringFlag(constants.FlagLabels, constants.FlagLabelsShort, map[string]string{}, "Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE")
	cmd.AddStringToStringFlag(constants.FlagAnnotations, constants.FlagAnnotationsShort, map[string]string{}, "Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE")

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request [seconds]")

	cmd.Command.SilenceUsage = true

	return cmd
}
