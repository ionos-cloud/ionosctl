package nodepool

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"
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
	cmd.AddCommand(NodepoolGetCmd())

	return cmd
}

func getNodepoolsPrint(c *core.CommandConfig, dcs *[]ionoscloud.NodePoolResponseData) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && dcs != nil {
		r.OutputJSON = dcs
		r.KeyValue = makeNodepoolPrintObj(dcs)                  // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
	}
	return r
}

type NodepoolPrint struct {
	Id                string `json:"Id,omitempty"`
	Name              string `json:"Name,omitempty"`
	NodeCount         int32  `json:"NodeCount,omitempty"`
	CoresCount        int32  `json:"CoresCount,omitempty"`
	CpuFamily         string `json:"CpuFamily,omitempty"`
	Ram               string `json:"Ram,omitempty"`
	AvailabilityZone  string `json:"AvailabilityZone,omitempty"`
	StorageSize       string `json:"StorageSize,omitempty"`
	MaintenanceWindow string `json:"MaintenanceWindow,omitempty"`
	Labels            string `json:"Labels,omitempty"`
	Annotations       string `json:"Annotations,omitempty"`
	State             string `json:"State,omitempty"`
}

var allCols = structs.Names(NodepoolPrint{})

func makeNodepoolPrintObj(clusters *[]ionoscloud.NodePoolResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*clusters))

	for _, cluster := range *clusters {
		var nodepoolPrint NodepoolPrint
		nodepoolPrint.Id = *cluster.GetId()

		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			nodepoolPrint.Name = *propertiesOk.GetName()
			nodepoolPrint.NodeCount = *propertiesOk.GetNodeCount()
			nodepoolPrint.CoresCount = *propertiesOk.GetCoresCount()
			nodepoolPrint.CpuFamily = *propertiesOk.GetCpuFamily()
			ramGb, err := utils.ConvertToGB(strconv.Itoa(int(*propertiesOk.GetRamSize())), utils.MegaBytes)
			if err == nil {
				nodepoolPrint.Ram = fmt.Sprintf("%d GB", ramGb)
			}
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
