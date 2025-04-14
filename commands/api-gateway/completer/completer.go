package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/spf13/viper"
	"strconv"
)

func GatewaysIDs() []string {

	if url := config.GetServerUrl(); url == constants.DefaultApiURL {
		viper.Set(constants.ArgServerUrl, "")
	}

	gateways, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysGet(context.Background()).Execute()
	if err != nil {
		return nil
	}

	Ids := []string{}
	for _, gateway := range gateways.GetItems() {
		gatewayId := gateway.Id
		Ids = append(Ids, gatewayId)
	}

	return Ids
}

func Routes(gatewayID string) []string {
	routesList, _, _ := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesGet(context.Background(), gatewayID).Execute()
	routesConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.ApiGatewayRoute, routesList)
	if err != nil {
		return nil
	}
	return completions.NewCompleter(routesConverted, "Id").AddInfo("Name").AddInfo("Paths").ToString()
}

func UpstreamsIDs(apigatewayId string, routeId string) []string {
	usedRoute, _, _ := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesFindById(context.Background(), apigatewayId, routeId).Execute()
	ids := []string{}
	for i := 0; i < len(usedRoute.Properties.Upstreams); i++ {
		ids = append(ids, strconv.Itoa(i))
	}
	return ids
}
