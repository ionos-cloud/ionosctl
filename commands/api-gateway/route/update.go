package route

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"

	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func RoutesPutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigatewayroute",
		Resource:  "route",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a route's properties. This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   "ionosctl apigateway route update --gateway-id GATEWAYID --route-id ROUTEID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagGatewayRouteID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			routeID := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayRouteID))
			r, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesFindById(context.Background(), apigatewayId, routeID).Execute()
			if err != nil {
				return fmt.Errorf("failed finding record: %w", err)
			}
			return partiallyUpdateGatewayAndPrint(c, r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescRoute, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			apigatewayId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.Routes(apigatewayId)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return ApiGatewayRouteCreateFlags(cmd)
}

func partiallyUpdateGatewayAndPrint(c *core.CommandConfig, r apigateway.RouteRead) error {
	input := r.Properties
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
		input.Name = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagType); viper.IsSet(fn) {
		input.Type = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPaths); viper.IsSet(fn) {
		input.Paths = viper.GetStringSlice(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMethods); viper.IsSet(fn) {
		input.Methods = viper.GetStringSlice(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagWebSocket); viper.IsSet(fn) {
		input.Websocket = pointer.From(viper.GetBool(fn))
	}
	if input.Upstreams == nil {
		input.Upstreams = make([]apigateway.RouteUpstreams, 1)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagScheme); viper.IsSet(fn) {
		input.Upstreams[0].Scheme = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagLoadBalancer); viper.IsSet(fn) {
		input.Upstreams[0].Loadbalancer = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
		input.Upstreams[0].Host = viper.GetString(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagPort); viper.IsSet(fn) {
		input.Upstreams[0].Port = viper.GetInt32(fn)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagWeight); viper.IsSet(fn) {
		input.Upstreams[0].Scheme = viper.GetString(fn)
	}

	apigatewayid := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
	routeid := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayRouteID))

	rn, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesPut(context.Background(),
		apigatewayid, routeid).RouteEnsure(apigateway.RouteEnsure{
		Properties: input,
		Id:         routeid,
	}).Execute()

	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApiGatewayRoute, rn,
		tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
