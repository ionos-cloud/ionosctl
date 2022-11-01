package cluster

import (
	"context"
	"github.com/spf13/viper"
	"os"

	"github.com/cjrd/allocate"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

var createProperties = ionoscloud.CreateClusterProperties{}
var conn = ionoscloud.Connection{}

func ClusterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create Mongo Clusters",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			var err error
			err = c.Command.Command.MarkFlagRequired(constants.FlagTemplateId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagInstances)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(constants.FlagLocation)
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
			err = c.Command.Command.MarkFlagRequired(constants.FlagMaintenanceTime)
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Creating Cluster...")
			day := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			time := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))

			maintenanceWindow := ionoscloud.MaintenanceWindow{}
			maintenanceWindow.SetDayOfTheWeek(ionoscloud.DayOfTheWeek(day))
			maintenanceWindow.SetTime(time)
			createProperties.SetMaintenanceWindow(maintenanceWindow)
			createProperties.SetConnections([]ionoscloud.Connection{conn})
			input := ionoscloud.CreateClusterRequest{}
			input.SetProperties(createProperties)

			cr, _, err := c.DbaasMongoServices.Clusters().Create(input)
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(c, &[]ionoscloud.ClusterResponse{cr}))
		},
		InitClient: true,
	})

	// Linked to properties struct
	_ = allocate.Zero(&createProperties)
	cmd.AddStringVarFlag(createProperties.DisplayName, constants.FlagName, constants.FlagNameShort, "", "The name of your cluster")
	cmd.AddStringVarFlag(createProperties.Location, constants.FlagLocation, constants.FlagLocationShort, "", "The physical location where the cluster will be created. This is the location where all your instances will be located. This property is immutable")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(createProperties.TemplateID, constants.FlagTemplateId, "", "", "The unique ID of the template, which specifies the number of cores, storage size, and memory")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoTemplateIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32VarFlag(createProperties.Instances, constants.FlagInstances, "", 0, "The total number of instances in the cluster (one primary and n-1 secondaries). Must be one of [3 5 7]")
	cmd.AddStringVarFlag(createProperties.MongoDBVersion, constants.FlagMongoVersion, "", "5.0", "The MongoDB version of your cluster.")

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	}) // TODO: Completions should be a flag option func

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Data Center creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Data Center creation [seconds]")

	// Connections
	_ = allocate.Zero(&conn)
	cmd.AddStringVarFlag(conn.DatacenterId, constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(conn.LanId, constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, cmd.Flag(constants.FlagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceVarFlag(conn.CidrList, constants.FlagCidrList, "", nil, "The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.64.0/18 10.233.0.0/18 10.233.114.0/24", core.RequiredFlagOption())
	/*
	   TemplateID        *string            `json:"templateID"`
	    MongoDBVersion    *string            `json:"mongoDBVersion,omitempty"`
	    Instances         *int32             `json:"instances"`
	    Connections       *[]Connection      `json:"connections"`
	    Location          *string            `json:"location"`
	    DisplayName       *string            `json:"displayName"`
	    MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"
	*/

	cmd.Command.SilenceUsage = true

	return cmd
}
