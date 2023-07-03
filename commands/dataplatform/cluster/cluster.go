package cluster

import (
	"fmt"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"
	"github.com/spf13/cobra"
)

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Dataplatform Cluster Operations",
			Long:             "This command allows you to interact with the already created clusters or creates new clusters in your virtual data center",
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	clusterCmd.Command.PersistentFlags().Bool(
		constants.ArgNoHeaders, false, "When using text output, don't print headers",
	)

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClustersKubeConfigCmd())

	return clusterCmd
}

func getClusterPrint(c *core.CommandConfig, dcs *[]ionoscloud.ClusterResponseData) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = getClusterRows(dcs)                        // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
	}
	return r
}

type ClusterPrint struct {
	Id                string `json:"Id,omitempty"`
	Name              string `json:"Name,omitempty"`
	Version           string `json:"Version,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	DatacenterId      string `json:"DatacenterId,omitempty"`
	State             string `json:"State,omitempty"`
}

var allCols = structs.Names(ClusterPrint{})

func getClusterRows(clusters *[]ionoscloud.ClusterResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))

	for _, cluster := range *clusters {
		var clusterPrint ClusterPrint
		clusterPrint.Id = *cluster.GetId()

		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			clusterPrint.Name = *propertiesOk.GetName()
			clusterPrint.DatacenterId = *propertiesOk.GetDatacenterId()
			clusterPrint.Version = *propertiesOk.GetDataPlatformVersion()
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
