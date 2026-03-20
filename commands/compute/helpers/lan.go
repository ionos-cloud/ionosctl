package helpers

import (
	"fmt"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func GetIPv6CidrBlockFromLAN(lan ionoscloud.Lan) (string, error) {
	if properties, ok := lan.GetPropertiesOk(); ok && properties != nil {
		if ipv6CidrBlock, ok := properties.GetIpv6CidrBlockOk(); ok && ipv6CidrBlock != nil {
			return *ipv6CidrBlock, nil
		} else if ok && ipv6CidrBlock == nil {
			return "", nil
		}
	}

	return "", fmt.Errorf("could not retrieve IPv6 Cidr Block from LAN")
}
