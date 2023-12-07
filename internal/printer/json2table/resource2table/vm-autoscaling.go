package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	ionoscloud "github.com/ionos-cloud/sdk-go-vm-autoscaling"
)

func ConvertVmAutoscalingGroupToTable(g ionoscloud.Group) ([]map[string]interface{}, error) {
	table, err := json2table.ConvertJSONToTable("", jsonpaths.VmAutoscalingGroup, g)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	ents, ok := g.GetEntitiesOk()
	if !ok || ents == nil {
		return table, nil
	}

	servers, ok := ents.Servers.GetItemsOk()
	if !ok || servers == nil {
		return table, nil
	}

	table[0]["Replicas"] = fmt.Sprintf("%d", len(*servers))

	return table, nil
}

func ConvertVmAutoscalingGroupsToTable(ls ionoscloud.GroupCollection) ([]map[string]interface{}, error) {
	items, ok := ls.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve items")
	}

	var clustersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertVmAutoscalingGroupToTable(item)
		if err != nil {
			return nil, err
		}

		clustersConverted = append(clustersConverted, temp...)
	}

	return clustersConverted, nil
}
