package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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

func Routes(gatewayID []string) []string {
	allRoutes := []string{}
	for _, gatewayID := range gatewayID {
		routesList, _, _ := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesGet(context.Background(), gatewayID).Execute()
		routesConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.ApiGatewayRoute, routesList)
		if err != nil {
			continue
		}
		allRoutes = append(allRoutes, completions.NewCompleter(routesConverted, "Id").AddInfo("Name").ToString()...)
	}
	return allRoutes
}
