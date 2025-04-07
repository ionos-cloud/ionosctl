package route

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/api-gateway/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/apigateway/v2"
	"github.com/spf13/viper"
)

func ApiGatewayRouteDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigatewayroute",
		Resource:  "route",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a gateway route",
		Example:   "ionosctl apigateway route delete --gateway-id ID --route-id ID_ROUTE",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagGatewayID, constants.FlagGatewayRouteID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			apigatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			routeId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayRouteID))
			z, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesFindById(context.Background(), apigatewayId, routeId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting route by id %s: %w", apigatewayId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete route %s ", z.Properties.Name),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().Apigateway.RoutesApi.ApigatewaysRoutesDelete(context.Background(), apigatewayId, routeId).Execute()
			if err != nil {
			}

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayRouteID, "", "", fmt.Sprintf("%s. Required or -%s", constants.DescRoute, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			print("hello 1 error")
			return completer.Routes(completer.GatewaysIDs())
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all routes. Required or -%s", constants.FlagRecordShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all routes!"))
	xs, _, err := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesGet(context.Background(), c.NS).Execute()

	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(z apigateway.RouteRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete routes %s ", z.Properties.Name),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().Apigateway.RoutesApi.ApigatewaysRoutesDelete(c.Context, c.NS, z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", z.Id, z.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
