package nodepool

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
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
	cmd.AddCommand(NodepoolCreateCmd())
	cmd.AddCommand(NodepoolGetCmd())
	cmd.AddCommand(NodepoolUpdateCmd())
	cmd.AddCommand(NodepoolDeleteCmd())

	return cmd
}

var (
	allJSONPaths = map[string]string{
		"Id":               "id",
		"Name":             "properties.name",
		"Nodes":            "",
		"Cores":            "properties.coresCount",
		"CpuFamily":        "properties.cpuFamily",
		"State":            "metadata.state",
		"AvailabilityZone": "properties.availabilityZone",
		"Labels":           "properties.labels",
		"Annotations":      "properties.annotations",
	}

	allCols = []string{"Id", "Name", "Nodes", "Cores", "CpuFamily", "Ram", "Storage", "MaintenanceWindow", "State",
		"AvailabilityZone", "Labels", "Annotations"}
	defaultCols = []string{"Id", "Name", "Nodes", "Cores", "CpuFamily", "Ram", "Storage", "MaintenanceWindow", "State"}
)

func convertNodePoolToTable(np ionoscloud.NodePoolResponseData) ([]map[string]interface{}, error) {
	properties, ok := np.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool properties")
	}

	ramRaw, ok := properties.GetRamSizeOk()
	if !ok || ramRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool RAM size")
	}

	gb := convbytes.Convert(int64(*ramRaw), convbytes.MB, convbytes.GB)
	ram := fmt.Sprintf("%v GB", gb)

	storageSizeRaw, ok := properties.GetStorageSizeOk()
	if !ok || storageSizeRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Storage size")
	}

	storageTypeRaw, ok := properties.GetStorageTypeOk()
	if !ok || storageTypeRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Storage type")
	}

	storageGb := convbytes.Convert(int64(*storageSizeRaw), convbytes.MB, convbytes.GB)
	storage := fmt.Sprintf("%v %v GB", *storageTypeRaw, storageGb)

	maintenanceWindowRaw, ok := properties.GetMaintenanceWindowOk()
	if !ok || maintenanceWindowRaw == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window")
	}

	day, ok := maintenanceWindowRaw.GetDayOfTheWeekOk()
	if !ok || day == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window Day")
	}

	time, ok := maintenanceWindowRaw.GetTimeOk()
	if !ok || time == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Maintenance Window Time")
	}

	maintenanceWindow := fmt.Sprintf("%v %v", *day, *time)

	temp, err := json2table.ConvertJSONToTable("", allJSONPaths, np)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["Ram"] = ram
	temp[0]["Storage"] = storage
	temp[0]["MaintenanceWindow"] = maintenanceWindow

	return temp, nil
}

func convertNodePoolsToTable(nps ionoscloud.NodePoolListResponseData) ([]map[string]interface{}, error) {
	var npsConverted = make([]map[string]interface{}, 0)

	items, ok := nps.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Node Pool Items")
	}

	for _, item := range *items {
		temp, err := convertNodePoolToTable(item)
		if err != nil {
			return nil, err
		}

		npsConverted = append(npsConverted, temp...)
	}

	return npsConverted, nil
}
