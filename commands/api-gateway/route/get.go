package route

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func RouteFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateatewayroute",
		Resource:  "route",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve routes",
		Example:   "ionosctl apigateway gateway get --gateway-id GATEWAYID --route ROUTE",
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
				return err
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.ApiGatewayRoute, r, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", constants.DescRoute, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			apigatewayId := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.Routes(apigatewayId)
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
