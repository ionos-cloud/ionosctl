package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"

	"github.com/cilium/fake"
	"github.com/cjrd/allocate"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	sdkdataplatform "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createProperties = sdkdataplatform.CreateClusterProperties{}

func ClusterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create Dataplatform Cluster",
		LongDesc:  "The cluster will be provisioned in the datacenter matching the provided datacenterID. Therefore the datacenter must be created upfront and must be accessible by the user issuing the request",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			err = c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagDatacenterId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagMaintenanceDay)
			if err != nil {
				return err
			}
			return c.Command.Command.MarkFlagRequired(constants.FlagMaintenanceTime)
		},
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Creating Cluster..."))

			day := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			time := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))

			maintenanceWindow := sdkdataplatform.MaintenanceWindow{}

			maintenanceWindow.SetDayOfTheWeek(day)
			maintenanceWindow.SetTime(time)
			createProperties.SetMaintenanceWindow(maintenanceWindow)

			input := sdkdataplatform.CreateClusterRequest{}
			input.SetProperties(createProperties)

			cr, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersPost(context.Background()).CreateClusterRequest(input).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, cr,
				printer.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Stdout, out)
			return nil
		},
		InitClient: true,
	})

	// Linked to properties struct
	_ = allocate.Zero(&createProperties)
	cmd.AddStringVarFlag(createProperties.Name, constants.FlagName, constants.FlagNameShort, "", "The name of your cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return fake.Names(10), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(createProperties.DataPlatformVersion, constants.FlagVersion, "", "23.7", "The version of your cluster")
	cmd.AddStringVarFlag(createProperties.DatacenterId, constants.FlagDatacenterId, constants.FlagIdShort, "", "The ID of the connected datacenter")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(nil), cobra.ShellCompDirectiveNoFileComp
	})

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
