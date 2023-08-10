package cluster

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var enumEditions = []string{"playground", "business", "enterprise"} // Remove whenever the SDK adds this as an actual type with enum vals
var enumTypes = []string{"replicaset", "sharded-cluster"}           // Remove whenever the SDK adds this as an actual type with enum vals

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
	clusterCmd.AddCommand(ClusterUpdateCmd())
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
	Health            string `json:"Health,omitempty"`
	State             string `json:"State,omitempty"`
	Instances         int32  `json:"Instances,omitempty"`
	MongoVersion      string `json:"MongoVersion,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	Location          string `json:"Location,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	LanId             string `json:"LanId,omitempty"`
	Cidr              string `json:"Cidr,omitempty"`
	TemplateId        string `json:"TemplateId,omitempty"`
	Cores             int32  `json:"Cores,omitempty"`
	RAM               string `json:"RAM,omitempty"`
	StorageSize       string `json:"StorageSize,omitempty"`
	StorageType       string `json:"StorageType,omitempty"`
}

var allCols = structs.Names(ClusterPrint{})

func getClusterRows(clusters *[]ionoscloud.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))

	for _, cluster := range *clusters {
		var clusterPrint ClusterPrint
		clusterPrint.ClusterId = *cluster.GetId()

		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			if propertiesOk.DisplayName != nil {
				clusterPrint.Name = *propertiesOk.GetDisplayName()
			}
			if propertiesOk.GetLocation() != nil {
				clusterPrint.Location = *propertiesOk.GetLocation()
			}
			if propertiesOk.GetTemplateID() != nil {
				clusterPrint.TemplateId = *propertiesOk.GetTemplateID()
			}
			if propertiesOk.GetConnectionString() != nil {
				clusterPrint.URL = *propertiesOk.GetConnectionString()
			}
			if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
				for _, vdcConnection := range *vdcConnectionsOk {
					if vdcConnection.GetDatacenterId() != nil {
						clusterPrint.DatacenterId = *vdcConnection.GetDatacenterId()
					}
					if vdcConnection.GetLanId() != nil {
						clusterPrint.LanId = *vdcConnection.GetLanId()
					}
					if vdcConnection.GetCidrList() != nil {
						clusterPrint.Cidr = strings.Join(*vdcConnection.GetCidrList(), ", ")
					}
				}
			}
			if propertiesOk.GetMongoDBVersion() != nil {
				clusterPrint.MongoVersion = *propertiesOk.GetMongoDBVersion()
			}
			if propertiesOk.GetInstances() != nil {
				clusterPrint.Instances = *propertiesOk.GetInstances()
			}
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				if maintenanceWindowOk.GetDayOfTheWeek() != nil && maintenanceWindowOk.GetTime() != nil {
					clusterPrint.MaintenanceWindow =
						fmt.Sprintf("%s %s", *maintenanceWindowOk.GetDayOfTheWeek(), *maintenanceWindowOk.GetTime())
				}
			}
			if f := propertiesOk.Ram; f != nil {
				clusterPrint.RAM = fmt.Sprintf("%d GB", convbytes.Convert(int64(*f), convbytes.MB, convbytes.GB))
			}
			if f := propertiesOk.StorageSize; f != nil {
				clusterPrint.StorageSize = fmt.Sprintf("%d GB", convbytes.Convert(int64(*f), convbytes.MB, convbytes.GB))
			}
			if f := propertiesOk.StorageType; f != nil {
				clusterPrint.StorageType = string(*f)
			}
			if f := propertiesOk.Cores; f != nil {
				clusterPrint.Cores = *f
			}
		}
		if md := cluster.Metadata; md != nil {
			if state := md.State; state != nil {
				clusterPrint.State = string(*state)
			}
			if health := md.Health; health != nil {
				clusterPrint.Health = string(*health)
			}
		}

		o := structs.Map(clusterPrint)
		out = append(out, o)
	}

	return out
}

func Clusters(fs ...Filter) (ionoscloud.ClusterList, error) {
	req := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		req = f(req)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return ionoscloud.ClusterList{}, fmt.Errorf("failed getting clusters: %w", err)
	}
	return clusters, err
}

type Filter func(ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
			req = req.Limit(viper.GetInt32(f))
		}
		if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
			req = req.Offset(viper.GetInt32(f))
		}
		return req
	}
}

func FilterNameFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
