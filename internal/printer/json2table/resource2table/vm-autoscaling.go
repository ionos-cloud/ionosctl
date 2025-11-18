package resource2table

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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

	var conv []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertVmAutoscalingGroupToTable(item)
		if err != nil {
			return nil, err
		}

		conv = append(conv, temp...)
	}

	return conv, nil
}

func ConvertVmAutoscalingServerToTable(sv ionoscloud.Server) ([]map[string]interface{}, error) {
	if sv.Properties == nil || sv.Properties.DatacenterServer == nil ||
		sv.Properties.DatacenterServer.Id == nil || sv.Properties.DatacenterServer.Href == nil {
		return nil, fmt.Errorf("server properties are incomplete: %+v", sv)
	}

	hrefFields := strings.FieldsFunc(*sv.Properties.DatacenterServer.Href, func(r rune) bool { return r == '/' })
	dcId := hrefFields[len(hrefFields)-3]
	cloudApiId := *sv.Properties.DatacenterServer.Id

	cloudApiServer, _, err := client.Must().CloudClient.ServersApi.DatacentersServersFindById(context.Background(), dcId, cloudApiId).Execute()
	if err != nil {
		return nil, fmt.Errorf("could not find server %s in datacenter %s via CloudAPI: %w", cloudApiId, dcId, err)
	}

	cloudApiServerAsTable, err := json2table.ConvertJSONToTable("", jsonpaths.Server, cloudApiServer)
	if err != nil {
		return nil, err
	}

	// Adding additional server info to each row
	cloudApiServerAsTable[0]["GroupServerId"] = *sv.Id
	return cloudApiServerAsTable, nil
}

// ConvertVmAutoscalingServersToTable converts a collection of servers to a table format.
func ConvertVmAutoscalingServersToTable(serverCollection ionoscloud.ServerCollection) ([]map[string]interface{}, error) {
	if serverCollection.Items == nil {
		return nil, fmt.Errorf("could not retrieve items")
	}

	s := *serverCollection.Items
	var table []map[string]interface{}

	for _, server := range s {
		serverTable, err := ConvertVmAutoscalingServerToTable(server)
		if err != nil {
			return nil, err
		}

		// Append each row of the server table to the main table
		table = append(table, serverTable...)
	}

	return table, nil
}
