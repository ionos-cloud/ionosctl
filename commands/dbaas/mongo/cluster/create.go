package cluster

import (
	"context"
	"os"

	"github.com/cjrd/allocate"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

const (
	flagName            = "name"
	flagNameShort       = "n"
	flagTemplateId      = "template-id"
	flagMongoVersion    = "mongo-version"
	flagInstances       = "instances"
	flagMaintenanceTime = "maintenance-time"
	flagMaintenanceDay  = "maintenance-day"
	flagLocation        = "location"
	flagLocationShort   = "l"
	flagDatacenterId    = "datacenter-id"
	flagCidrList        = "cidr-list"
	flagLanId           = "lan-id"
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
			err = c.Command.Command.MarkFlagRequired(flagTemplateId)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(flagName)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(flagInstances)
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired(flagLocation)
			if err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Creating Cluster...")
			day, err := c.Command.Command.Flags().GetString(flagMaintenanceDay)
			if err != nil {
				return err
			}
			time, err := c.Command.Command.Flags().GetString(flagMaintenanceTime)
			if err != nil {
				return err
			}

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
	cmd.AddStringVarFlag(createProperties.DisplayName, flagName, flagNameShort, "", "The name of your cluster")
	cmd.AddStringVarFlag(createProperties.Location, flagLocation, flagLocationShort, "", "The physical location where the cluster will be created. This is the location where all your instances will be located. This property is immutable")
	_ = cmd.Command.RegisterFlagCompletionFunc(flagLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(createProperties.TemplateID, flagTemplateId, "", "", "The unique ID of the template, which specifies the number of cores, storage size, and memory")
	_ = cmd.Command.RegisterFlagCompletionFunc(flagTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoTemplateIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32VarFlag(createProperties.Instances, flagInstances, "", 0, "The total number of instances in the cluster (one primary and n-1 secondaries). Must be one of [3 5 7]")
	cmd.AddStringVarFlag(createProperties.MongoDBVersion, flagMongoVersion, "", "5.0", "The MongoDB version of your cluster.")

	// Maintenance
	cmd.AddStringFlag(flagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	cmd.AddStringFlag(flagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	}) // TODO: Completions should be a flag option func

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Data Center creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Data Center creation [seconds]")

	// Connections
	_ = allocate.Zero(&conn)
	cmd.AddStringVarFlag(conn.DatacenterId, flagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringVarFlag(conn.LanId, flagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(flagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, cmd.Flag(flagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceVarFlag(conn.CidrList, flagCidrList, "", nil, "The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.64.0/18 10.233.0.0/18 10.233.114.0/24", core.RequiredFlagOption())
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
