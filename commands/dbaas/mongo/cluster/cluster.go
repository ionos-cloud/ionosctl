package cluster

import (
	"fmt"
	"strings"

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
	//clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
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
		r.KeyValue = getClusterRows(dcs)                                                                                   // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}

type ClusterPrint struct {
	ClusterId         string `json:"ClusterId,omitempty"`
	TemplateId        string `json:"TemplateId,omitempty"`
	DisplayName       string `json:"DisplayName,omitempty"`
	URL               string `json:"URL,omitempty"`
	State             string `json:"State,omitempty"`
	Instances         int32  `json:"Instances,omitempty"`
	Location          string `json:"Location,omitempty"`
	MongoVersion      string `json:"MongoVersion,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	LanId             string `json:"LanId,omitempty"`
	CidrList          string `json:"CidrList,omitempty"`
}

var allCols = structs.Names(ClusterPrint{})

func getClusterRows(clusters *[]ionoscloud.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))
	for _, cluster := range *clusters {
		var clusterPrint ClusterPrint
		clusterPrint.ClusterId = *cluster.GetId()
		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			clusterPrint.DisplayName = *propertiesOk.GetDisplayName()
			clusterPrint.Location = *propertiesOk.GetLocation()
			clusterPrint.TemplateId = *propertiesOk.GetTemplateID()
			clusterPrint.URL = *propertiesOk.GetConnectionString()
			if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
				for _, vdcConnection := range *vdcConnectionsOk {
					// TODO: This only gets the last items in the connections slice. DBaaS API seems to only support one connection atm.
					// Create connections sub-command if multiple connections are allowed
					clusterPrint.DatacenterId = *vdcConnection.GetDatacenterId()
					clusterPrint.LanId = *vdcConnection.GetLanId()
					clusterPrint.CidrList = strings.Join(*vdcConnection.GetCidrList(), ", ")
				}
			}
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
