package resource2table

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
)

func ConvertDNSSECToTable(keys dns.DnssecKeyReadList) ([]map[string]interface{}, error) {
	table, err := json2table.ConvertJSONToTable("", jsonpaths.DnsSecKey, keys)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	if keys.Metadata == nil || keys.Metadata.Items == nil || len(keys.Metadata.Items) == 0 {
		return table, nil
	}

	for i, item := range keys.Metadata.Items {
		table[i]["Id"] = *keys.Id
		table[i]["KeyTag"] = *item.KeyTag
		table[i]["DigestAlgorithmMnemonic"] = *item.DigestAlgorithmMnemonic
		table[i]["Digest"] = *item.Digest
		table[i]["Flags"] = *item.KeyData.Flags
		table[i]["PubKey"] = *item.KeyData.PubKey
		table[i]["ComposedKeyData"] = *item.ComposedKeyData
		table[i]["Algorithm"] = *keys.Properties.KeyParameters.Algorithm
		table[i]["NsecMode"] = *keys.Properties.NsecParameters.NsecMode
	}

	return table, nil
}
