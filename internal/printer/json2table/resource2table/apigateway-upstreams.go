package resource2table

import (
	"fmt"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
)

func ConverApiGatewayUpstreamsToTable(upstream []apigateway.RouteUpstreams) []map[string]interface{} {
	if len(upstream) == 0 {
		return nil // empty output
	}

	var convertedUpstreams = make([]map[string]interface{}, len(upstream))
	for idx, ups := range upstream {
		convertedUpstreams[idx] = make(map[string]interface{})

		convertedUpstreams[idx]["UpstreamId"] = fmt.Sprintf("%v", idx)
		convertedUpstreams[idx]["Scheme"] = ups.GetScheme()
		convertedUpstreams[idx]["Loadbalancer"] = ups.GetLoadbalancer()
		convertedUpstreams[idx]["Host"] = ups.GetHost()
		convertedUpstreams[idx]["Port"] = ups.GetPort()
		convertedUpstreams[idx]["Weight"] = ups.GetWeight()
	}

	return convertedUpstreams
}

func ConvertApiGatewayCustomDomainsToTable(customDomains []apigateway.GatewayCustomDomains) []map[string]interface{} {
	if len(customDomains) == 0 {
		return nil
	}

	var convertedCustomDomains = make([]map[string]interface{}, len(customDomains))
	for idx, customDomain := range customDomains {
		convertedCustomDomains[idx] = make(map[string]interface{})

		convertedCustomDomains[idx]["CustomDomainsId"] = fmt.Sprintf("%v", idx)
		convertedCustomDomains[idx]["Name"] = customDomain.GetName()
		convertedCustomDomains[idx]["CertificateId"] = customDomain.GetCertificateId()
	}
	return convertedCustomDomains
}
