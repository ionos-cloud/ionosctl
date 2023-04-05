package cluster

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo"
	"github.com/spf13/cobra"
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Mongo Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas mongo cluster` allow you to manage the Mongo Clusters under your account.",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	clusterCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd()) // TODO
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

// TODO: should be moved to printer package as a decoupled func, to reduce duplication
func getClusterPrint(c *core.CommandConfig, dcs *[]ionoscloud.ClusterResponse) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = getClusterRows(dcs)                            // map header -> rows
		r.Columns = printer.GetHeaders(allCols, allCols[0:6], cols) // headers
	}
	return r
}

type ClusterPrint struct {
	ClusterId         string `json:"ClusterId,omitempty"`
	Name              string `json:"Name,omitempty"`
	URL               string `json:"URL,omitempty"`
	State             string `json:"State,omitempty"`
	Instances         int32  `json:"Instances,omitempty"`
	MongoVersion      string `json:"MongoVersion,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	Location          string `json:"Location,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	LanId             string `json:"LanId,omitempty"`
	Cidr              string `json:"Cidr,omitempty"`
	TemplateId        string `json:"TemplateId,omitempty"`
}

var allCols = structs.Names(ClusterPrint{})

func getClusterRows(clusters *[]ionoscloud.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))

	for _, cluster := range *clusters {
		var clusterPrint ClusterPrint
		clusterPrint.ClusterId = *cluster.GetId()

		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			clusterPrint.Name = *propertiesOk.GetDisplayName()
			clusterPrint.Location = *propertiesOk.GetLocation()
			clusterPrint.TemplateId = *propertiesOk.GetTemplateID()
			clusterPrint.URL = *propertiesOk.GetConnectionString()
			if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
				for _, vdcConnection := range *vdcConnectionsOk {
					clusterPrint.DatacenterId = *vdcConnection.GetDatacenterId()
					clusterPrint.LanId = *vdcConnection.GetLanId()
					clusterPrint.Cidr = strings.Join(*vdcConnection.GetCidrList(), ", ")
				}
			}
			clusterPrint.MongoVersion = *propertiesOk.GetMongoDBVersion()
			clusterPrint.Instances = *propertiesOk.GetInstances()
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				clusterPrint.MaintenanceWindow =
					fmt.Sprintf("%s %s", *maintenanceWindowOk.GetDayOfTheWeek(), *maintenanceWindowOk.GetTime())
			}
		}
		clusterPrint.State = string(*cluster.GetMetadata().GetState())
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
