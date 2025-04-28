package upstreams

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"
)

func AddCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "upstreams",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Upstreams consist of schme, loadbalancer, host, port and weight",
		Example:   "ionosctl apigateway route upstreams add --gateway-id ID --route-id ID_ROUTE --host HOST",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagGatewayRouteID, constants.FlagHost); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			routeId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayRouteID))
			usedRoute, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesFindById(context.Background(), apigatewayId, routeId).Execute()
			input := usedRoute.Properties
			elem := len(input.Upstreams)

			if input.Upstreams == nil {
				input.Upstreams = make([]apigateway.RouteUpstreams, 1)
			}
			input.Upstreams = append(input.Upstreams, apigateway.RouteUpstreams{})
			elem = len(input.Upstreams) - 1

			if fn := core.GetFlagName(c.NS, constants.FlagScheme); true {
				input.Upstreams[elem].Scheme = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLoadBalancer); true {
				input.Upstreams[elem].Loadbalancer = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
				input.Upstreams[elem].Host = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagPort); true {
				input.Upstreams[elem].Port = viper.GetInt32(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagWeight); viper.IsSet(fn) {
				input.Upstreams[elem].Scheme = viper.GetString(fn)
			}
			input.Name = usedRoute.Properties.Name
			input.Websocket = usedRoute.Properties.Websocket
			input.Paths = usedRoute.Properties.Paths
			input.Type = usedRoute.Properties.Type
			input.Methods = usedRoute.Properties.Methods

			rec, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesPut(context.Background(), apigatewayId, routeId).
				RouteEnsure(apigateway.RouteEnsure{
					Id:         routeId,
					Properties: input,
				}).Execute()

			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			upstreamsConverted := resource2table.ConverApiGatewayUpstreamsToTable(rec.Properties.Upstreams)

			out, err := jsontabwriter.GenerateOutputPreconverted(rec.Properties.Upstreams, upstreamsConverted,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return ApiGatewayRouteCreateFlags(cmd)
}
func ApiGatewayRouteCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescRoute, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			apigatewayId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.Routes(apigatewayId)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)
	cmd.AddSetFlag(constants.FlagScheme, "s", "http",
		[]string{"http", "https", "grpc", "grpcs"},
		"The target URL of the upstream.")
	cmd.AddStringFlag(constants.FlagLoadBalancer, "", "roundrobin", "The load balancer algorithm.")
	cmd.AddStringFlag(constants.FlagHost, "", "", "The host of the upstream. Field is validated as hostname according to RFC1123.", core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagPort, "", 80, "The port of the upstream.")
	cmd.AddInt32Flag(constants.FlagWeight, "", 100, "Weight with which to split traffic to the upstream.")
	return cmd
}
