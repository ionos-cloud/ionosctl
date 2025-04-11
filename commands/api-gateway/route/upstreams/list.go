package upstreams

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func ListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "upstreams",
		Verb:      "list",
		Aliases:   []string{"l"},
		ShortDesc: "Upstreams consist of schme, loadbalancer, host, port and weight",
		Example:   "ionosctl apigateway route upstreams list --gateway-id ID --route-id ID_ROUTE",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagGatewayRouteID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			apiGatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			g, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesGet(context.Background(), apiGatewayId).Execute()
			if err != nil {
				return err
			}
			g.GetItems()
			//apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			//g, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apigatewayId).Execute()
			//if err != nil {
			//	return err
			//}
			return nil //partiallyUpdateGatewayPrint(c, g)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagGatewayShort, "", constants.DescGateway, core.RequiredFlagOption(),
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

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
