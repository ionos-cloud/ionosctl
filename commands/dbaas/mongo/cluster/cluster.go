package cluster

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterDeleteAllCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}

// TODO: Why is this tightly coupled to resources.ClusterResponse? Should just take Headers and Columns as params. should also be moved to printer package, to reduce duplication
//
// this is a nightmare to maintain if it is tightly coupled to every single resource!!!!!!!!!!!!
func getClusterPrint(c *core.CommandConfig, dcs *[]ionoscloud.ClusterResponse) printer.Result {
	r := printer.Result{}
	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = getClusterRows(dcs)                                                                                       // map header -> rows
		r.Columns = printer.GetHeaders(allCols, allCols[0:6], viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}

type ClusterPrint struct {
	ClusterId         string `json:"ClusterId,omitempty"`
	Location          string `json:"Location,omitempty"`
	TemplateId        string `json:"TemplateId,omitempty"`
	State             string `json:"State,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty"`
	MongoVersion      string `json:"MongoVersion,omitempty"`
	Instances         int32  `json:"Instances,omitempty"`
	Connections       string `json:"Connections,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
}

var allCols = structs.Names(ClusterPrint{})

func getClusterRows(clusters *[]ionoscloud.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))
	for _, cluster := range *clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetDisplayNameOk(); ok && displayNameOk != nil {
				clusterPrint.DisplayName = *displayNameOk
			}
			if locationOk, ok := propertiesOk.GetLocationOk(); ok && locationOk != nil {
				clusterPrint.Location = string(*locationOk)
			}
			if templateIdOk, ok := propertiesOk.GetTemplateIDOk(); ok && templateIdOk != nil {
				clusterPrint.TemplateId = string(*templateIdOk)
			}
			if connectionsOk, ok := propertiesOk.GetConnectionStringOk(); ok && connectionsOk != nil {
				clusterPrint.Connections = *connectionsOk
			}
			//if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
			//	for _, vdcConnection := range *vdcConnectionsOk {
			//		// TODO: This seems to only get the last items in the connections slice?
			//		if vdcIdOk, ok := vdcConnection.GetDatacenterIdOk(); ok && vdcIdOk != nil {
			//			clusterPrint.DatacenterId = *vdcIdOk
			//		}
			//		if lanIdOk, ok := vdcConnection.GetLanIdOk(); ok && lanIdOk != nil {
			//			clusterPrint.LanId = *lanIdOk
			//		}
			//		if ipAddressOk, ok := vdcConnection.GetCidrOk(); ok && ipAddressOk != nil {
			//			clusterPrint.Cidr = *ipAddressOk
			//		}
			//	}
			//}
			if versionOk, ok := propertiesOk.GetMongoDBVersionOk(); ok && versionOk != nil {
				clusterPrint.MongoVersion = *versionOk
			}
			if replicasOk, ok := propertiesOk.GetInstancesOk(); ok && replicasOk != nil {
				clusterPrint.Instances = *replicasOk
			}
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				var maintenanceWindow string
				if weekdayOk, ok := maintenanceWindowOk.GetDayOfTheWeekOk(); ok && weekdayOk != nil {
					maintenanceWindow = string(*weekdayOk)
				}
				if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
					maintenanceWindow = fmt.Sprintf("%s %s", maintenanceWindow, *timeOk)
				}
				clusterPrint.MaintenanceWindow = maintenanceWindow
			}
		}
		if metadataOk, ok := cluster.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				clusterPrint.State = string(*stateOk)
			}
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
