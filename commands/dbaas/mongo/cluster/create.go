package cluster

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
)

var clusterProperties = ionoscloud.CreateClusterProperties{} // flag values will point to these variables in memory

func ClusterCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil /* circular dependency 🤡*/, core.CommandBuilder{
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create Mongo Clusters",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Creating Cluster...")
			input := ionoscloud.CreateClusterRequest{}

			cr, r, err := c.DbaasMongoServices.Clusters().Create(input)
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(r, c, []CreateRes))
		},
		InitClient: true,
	})

	// TODO: Move ArgName to DBAAS level constants
	cmd.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddStringFlag(dbaaspg.ArgLocation, "", "", "Location")
	cmd.AddStringFlag("template-id", "", "", "Template to use")
	cmd.AddInt32Flag("instances", "", 0, "Instances")
	cmd.AddStringFlag(dbaaspg.ArgMaintenanceTime, dbaaspg.ArgMaintenanceTimeShort, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59")
	cmd.AddStringFlag(dbaaspg.ArgMaintenanceDay, dbaaspg.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = cmd.Command.RegisterFlagCompletionFunc(dbaaspg.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Data Center creation to be executed")
	cmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Data Center creation [seconds]")
	cmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
	   TemplateID        *string            `json:"templateID"`
	    MongoDBVersion    *string            `json:"mongoDBVersion,omitempty"`
	    Instances         *int32             `json:"instances"`
	    Connections       *[]Connection      `json:"connections"`
	    Location          *string            `json:"location"`
	    DisplayName       *string            `json:"displayName"`
	    MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"
	*/

	return cmd
}
