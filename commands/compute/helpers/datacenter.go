package helpers

import (
	"fmt"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func GetDataCenters(datacenters resources.Datacenters) []resources.Datacenter {
	dc := make([]resources.Datacenter, 0)
	if items, ok := datacenters.GetItemsOk(); ok && items != nil {
		for _, datacenter := range *items {
			dc = append(dc, resources.Datacenter{Datacenter: datacenter})
		}
	}
	return dc
}

func GetIPv6CidrBlockFromDatacenter(dc ionoscloud.Datacenter) (string, error) {
	if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
		if ipv6CidrBlock, ok := properties.GetIpv6CidrBlockOk(); ok && ipv6CidrBlock != nil {
			return *ipv6CidrBlock, nil
		}
	}

	return "", fmt.Errorf("could not get IPv6 Cidr Block from Datacenter")
}

func GetGroups(groups resources.Groups) []resources.Group {
	u := make([]resources.Group, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.Group{Group: item})
		}
	}
	return u
}
