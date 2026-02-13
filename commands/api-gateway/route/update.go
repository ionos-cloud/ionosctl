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

// Note: viper is still used in completion function on line 50

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
			apigatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}

			routeID, err := c.Command.Command.Flags().GetString(constants.FlagGatewayRouteID)
			if err != nil {
				return err
			}

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
	if c.Command.Command.Flags().Changed(constants.FlagName) {
		name, err := c.Command.Command.Flags().GetString(constants.FlagName)
		if err != nil {
			return err
		}
		input.Name = name
	}
	if c.Command.Command.Flags().Changed(constants.FlagType) {
		routeType, err := c.Command.Command.Flags().GetString(constants.FlagType)
		if err != nil {
			return err
		}
		input.Type = routeType
	}
	if c.Command.Command.Flags().Changed(constants.FlagPaths) {
		paths, err := c.Command.Command.Flags().GetStringSlice(constants.FlagPaths)
		if err != nil {
			return err
		}
		input.Paths = paths
	}
	if c.Command.Command.Flags().Changed(constants.FlagMethods) {
		methods, err := c.Command.Command.Flags().GetStringSlice(constants.FlagMethods)
		if err != nil {
			return err
		}
		input.Methods = methods
	}
	if c.Command.Command.Flags().Changed(constants.FlagWebSocket) {
		websocket, err := c.Command.Command.Flags().GetBool(constants.FlagWebSocket)
		if err != nil {
			return err
		}
		input.Websocket = pointer.From(websocket)
	}
	if input.Upstreams == nil {
		input.Upstreams = make([]apigateway.RouteUpstreams, 1)
	}
	if c.Command.Command.Flags().Changed(constants.FlagScheme) {
		scheme, err := c.Command.Command.Flags().GetString(constants.FlagScheme)
		if err != nil {
			return err
		}
		input.Upstreams[0].Scheme = scheme
	}
	if c.Command.Command.Flags().Changed(constants.FlagLoadBalancer) {
		loadbalancer, err := c.Command.Command.Flags().GetString(constants.FlagLoadBalancer)
		if err != nil {
			return err
		}
		input.Upstreams[0].Loadbalancer = loadbalancer
	}
	if c.Command.Command.Flags().Changed(constants.FlagHost) {
		host, err := c.Command.Command.Flags().GetString(constants.FlagHost)
		if err != nil {
			return err
		}
		input.Upstreams[0].Host = host
	}
	if c.Command.Command.Flags().Changed(constants.FlagPort) {
		port, err := c.Command.Command.Flags().GetInt32(constants.FlagPort)
		if err != nil {
			return err
		}
		input.Upstreams[0].Port = port
	}
	if c.Command.Command.Flags().Changed(constants.FlagWeight) {
		weight, err := c.Command.Command.Flags().GetString(constants.FlagWeight)
		if err != nil {
			return err
		}
		input.Upstreams[0].Scheme = weight
	}

	apigatewayid, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
	if err != nil {
		return err
	}

	routeid, err := c.Command.Command.Flags().GetString(constants.FlagGatewayRouteID)
	if err != nil {
		return err
	}

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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
