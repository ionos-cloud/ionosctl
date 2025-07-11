package upstreams

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"
)

func RemovetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "upstreams",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Upstreams consist of schme, loadbalancer, host, port and weight",
		Example:   "ionosctl apigateway route upstreams remove --gateway-id ID --route-id ID_ROUTE --upstream-id UPSTREAMID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagGatewayRouteID, constants.FlagUpstreamId); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			routeId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayRouteID))
			upstreamId := viper.GetInt(core.GetFlagName(c.NS, constants.FlagUpstreamId))
			usedRoutes, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesFindById(context.Background(), apigatewayId, routeId).Execute()
			if err != nil {
				return err
			}
			input := usedRoutes.Properties
			if input.Upstreams == nil || len(input.Upstreams) == 0 {
				return fmt.Errorf("There are no upstreams defined in this route!")
			}
			if upstreamId < 0 || upstreamId >= len(input.Upstreams) {
				return fmt.Errorf("Invalid Upstreams index")
			}
			input.Upstreams = append(input.Upstreams[:upstreamId], input.Upstreams[upstreamId+1:]...)

			_, _, err = client.Must().Apigateway.RoutesApi.ApigatewaysRoutesPut(context.Background(), apigatewayId, routeId).
				RouteEnsure(apigateway.RouteEnsure{
					Id:         routeId,
					Properties: input,
				}).Execute()
			if err != nil {
				return err
			}
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescRoute, constants.FlagAllShort), core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			apigatewayId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.Routes(apigatewayId)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagUpstreamId, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescUpstream, constants.FlagAllShort), core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			apigatewayId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			routeId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayRouteID))
			return completer.UpstreamsIDs(apigatewayId, routeId)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations))

	return cmd
}
