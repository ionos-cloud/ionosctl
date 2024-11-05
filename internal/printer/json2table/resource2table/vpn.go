package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
)

func ConvertVPNWireguardGatewaysToTable(gateways vpn.WireguardGatewayReadList) ([]map[string]interface{}, error) {
	items, ok := gateways.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Users items")
	}

	var usersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := ConvertVPNWireguardGatewayToTable(item)
		if err != nil {
			return nil, err
		}

		usersConverted = append(usersConverted, temp...)
	}

	return usersConverted, nil
}

func ConvertVPNWireguardGatewayToTable(gateway vpn.WireguardGatewayRead) ([]map[string]interface{}, error) {
	table, err := json2table.ConvertJSONToTable("", jsonpaths.VPNWireguardGateway, gateway)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	for _, connection := range *gateway.Properties.Connections {
		table[0]["DatacenterId"] = *connection.DatacenterId
		table[0]["LanId"] = *connection.LanId
		table[0]["ConnectionIPv4"] = *connection.Ipv4CIDR
		table[0]["ConnectionIPv6"] = *connection.Ipv6CIDR
	}

	return table, nil
}
