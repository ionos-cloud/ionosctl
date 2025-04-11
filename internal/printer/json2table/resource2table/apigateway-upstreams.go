package resource2table

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
)

func ConvertApiGatewayUpstreamToTable(ups apigateway.RouteRead) ([]map[string]interface{}, error) {
	properties, ok := ups.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve ApiGateway Properties")
	}

	upstr, ok := properties.GetUpstreamsOk()
	if !ok || upstr == nil {
		return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams")
	}
	var upstreamsFormatted []interface{}
	for _, upstream := range upstr {
		scheme, ok := upstream.GetSchemeOk()
		if !ok || scheme == nil {
			return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams name")
		}

		loadbalancer, ok := upstream.GetLoadbalancerOk()
		if !ok || loadbalancer == nil {
			return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams loadbalancer")
		}

		host, ok := upstream.GetHostOk()
		if !ok || host == nil {
			return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams host")
		}

		port, ok := upstream.GetPortOk()
		if !ok || port == nil {
			return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams port")
		}

		weight, ok := upstream.GetWeightOk()
		if !ok || weight == nil {
			return nil, fmt.Errorf("could not retrieve ApiGateway Upstreams weight")
		}
		upstreamsFormatted = append(upstreamsFormatted, fmt.Sprintf("%v (%v) (%v) (%v) (%v)", *scheme, *loadbalancer, *host, *port, *weight))
	}

	propertiesUpstream, upOk := properties.GetUpstreamsOk()
	if !upOk || propertiesUpstream == nil {
		convertedUpstreams[0]["Upstreams"] = upstreamsFormatted
	}

	return convertedUpstreams, nil
}
