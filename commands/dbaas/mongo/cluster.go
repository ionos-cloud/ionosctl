package mongo

import (
	"context"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/services/dbaas-mongo/resources"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Mongo Clusters",
		LongDesc:  "Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:   "ionosctl dbaas mongo cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Clusters...")
			clusters, _, err := c.DbaasMongoServices.Clusters().List("")
			if err != nil {
				return err
			}

			return c.Printer.Print(getClusterPrint(nil, c, clusters.GetItems()))

		},
		InitClient: true,
	})
	// TODO: Move ArgName to DBAAS level constants
	list.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	list.AddStringSliceFlag(config.ArgCols, "", allCols[0:6], printer.ColsMessage(allCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	return clusterCmd
}

// TODO: Why is this tightly coupled to resources.ClusterResponse? Should just take Headers and Columns as params
// TODO: should also be moved to printer package, to reduce duplication
func getClusterPrint(resp *resources.Response, c *core.CommandConfig, dcs *[]ionoscloud.ClusterResponse) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getClusterRows(dcs)                                                            // map header -> rows
			r.Columns = getClusterHeaders(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgCols))) // headers
		}
	}
	return r
}

var allCols = structs.Names(ClusterPrint{})

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

func getClusterHeaders(customColumns []string) []string {
	if customColumns == nil {
		return allCols[0:6]
	}
	//for _, c := customColumns {
	//	if slices.Contains(allCols, c) {}
	//}
	return customColumns
}
