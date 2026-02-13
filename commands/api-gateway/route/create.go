package route

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ApiGatewayRoutesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigatewayroute",
		Resource:  "route",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a new route",
		Example:   "ionosctl apigateway route create --gateway-id ID --name NAME --paths PATHS --methods METHODS --host HOST",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagName, constants.FlagPaths, constants.FlagMethods, constants.FlagHost); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := apigateway.Route{}

			if err := modifyRoutePropertiesFromFlags(c, &input); err != nil {
				return err
			}

			apigatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}

			rec, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesPost(context.Background(), apigatewayId).
				RouteCreate(apigateway.RouteCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApiGatewayRoute, rec,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
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
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the route.", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagType, "", "http", " Default: http. This field specifies the protocol used by the ingress to route traffic to the backend service.")
	cmd.AddStringFlag(constants.FlagPaths, "", "", fmt.Sprintf("The paths that the route should match."), core.RequiredFlagOption())
	cmd.AddStringSliceFlag(constants.FlagMethods, "m", []string{},
		"The HTTP methods that the route should match.", core.RequiredFlagOption(), core.WithCompletion(
			func() []string {
				return []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD", "CONNECT", "TRACE"}
			}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)
	cmd.AddBoolFlag(constants.FlagWebSocket, "", false, "To enable websocket support.")
	cmd.AddSetFlag(constants.FlagScheme, "s", "http",
		[]string{"http", "https", "grpc", "grpcs"},
		"The target URL of the upstream.", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagLoadBalancer, "", "roundrobin", "The load balancer algorithm.")
	cmd.AddStringFlag(constants.FlagHost, "", "", "The host of the upstream. Field is validated as hostname according to RFC1123.", core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagPort, "", 80, "The port of the upstream.")
	cmd.AddInt32Flag(constants.FlagWeight, "", 100, "Weight with which to split traffic to the upstream.")
	return cmd
}

func modifyRoutePropertiesFromFlags(c *core.CommandConfig, input *apigateway.Route) error {
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
	return nil
}
