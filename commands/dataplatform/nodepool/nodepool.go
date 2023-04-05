package nodepool

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
	"github.com/spf13/cobra"
)

func NodepoolCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "nodepool",
			Aliases:          []string{"np"},
			Short:            "Dataplatform Nodepool Operations",
			Long:             "Node pools are the resources that powers the DataPlatformCluster",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	cmd.AddCommand(NodepoolListCmd())
	cmd.AddCommand(NodepoolCreateCmd())
	cmd.AddCommand(NodepoolGetCmd())
	cmd.AddCommand(NodepoolUpdateCmd())
	cmd.AddCommand(NodepoolDeleteCmd())

	return cmd
}

func getNodepoolsPrint(c *core.CommandConfig, dcs *[]ionoscloud.NodePoolResponseData) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = makeNodepoolPrintObj(dcs)                 // map header -> rows
		r.Columns = printer.GetHeaders(allCols, defCols, cols) // headers
	}
	return r
}

type NodepoolPrint struct {
	Id                string `json:"Id,omitempty"`
	Name              string `json:"Name,omitempty"`
	Nodes             int32  `json:"Nodes,omitempty"`
	Cores             int32  `json:"Cores,omitempty"`
	CpuFamily         string `json:"CpuFamily,omitempty"`
	Ram               string `json:"Ram,omitempty"`
	Storage           string `json:"Storage,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	State             string `json:"State,omitempty"`

	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	Labels           string `json:"Labels,omitempty"`
	Annotations      string `json:"Annotations,omitempty"`
}

var allCols = structs.Names(NodepoolPrint{})
var defCols = allCols[:len(allCols)-3]

func makeNodepoolPrintObj(clusters *[]ionoscloud.NodePoolResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))

	for _, cluster := range *clusters {
		var nodepoolPrint NodepoolPrint
		nodepoolPrint.Id = *cluster.GetId()

		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			nodepoolPrint.Name = *propertiesOk.GetName()
			nodepoolPrint.Nodes = *propertiesOk.GetNodeCount()
			nodepoolPrint.Cores = *propertiesOk.GetCoresCount()
			nodepoolPrint.CpuFamily = *propertiesOk.GetCpuFamily()
			gb, err := utils.ConvertToGB(strconv.Itoa(int(*propertiesOk.GetRamSize())), utils.MegaBytes)
			if err == nil {
				nodepoolPrint.Ram = fmt.Sprintf("%d GB", gb)
			}
			nodepoolPrint.Storage = fmt.Sprintf("%s %d GB", *propertiesOk.StorageType, *propertiesOk.GetStorageSize())
			nodepoolPrint.AvailabilityZone = *propertiesOk.GetCpuFamily()
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				nodepoolPrint.MaintenanceWindow =
					fmt.Sprintf("%s %s", *maintenanceWindowOk.GetDayOfTheWeek(), *maintenanceWindowOk.GetTime())
			}
			nodepoolPrint.Labels = ""
			for k, v := range *propertiesOk.GetLabels() {
				nodepoolPrint.Labels += fmt.Sprintf("%s = %s; ", k, v)
			}
			nodepoolPrint.Annotations = ""
			for k, v := range *propertiesOk.GetAnnotations() {
				nodepoolPrint.Annotations += fmt.Sprintf("%s = %s; ", k, v)
			}

		}
		nodepoolPrint.State = string(*cluster.GetMetadata().GetState())
		o := structs.Map(nodepoolPrint)
		out = append(out, o)
	}
	return out
}
