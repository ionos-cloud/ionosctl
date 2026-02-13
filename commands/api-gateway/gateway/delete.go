package gateway

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

// Note: viper is still used for global flag constants.ArgForce

func ApiGatewayDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "apigateway",
		Resource:  "gateway",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a gateway",
		Example:   "ionosctl apigateway gateway delete --gateway-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagGatewayID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			all, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
			if err != nil {
				return err
			}

			if all {
				return deleteAll(c)
			}

			apigatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}

			z, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysFindById(context.Background(), apigatewayId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting gateway by id %s: %w", apigatewayId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete gateway with name: %s, id: %s ", z.Properties.Name, z.Id),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().Apigateway.APIGatewaysApi.ApigatewaysDelete(context.Background(), apigatewayId).Execute()
			if err != nil {
				return err
			}

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", fmt.Sprintf("%s. Required or -%s", constants.DescGateway, constants.ArgAllShort),
		core.WithCompletion(func() []string {
			return completer.GatewaysIDs()
		}, constants.ApiGatewayRegionalURL, constants.GatewayLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all gateways. Required or -%s", constants.FlagGatewayShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all gateways!"))
	xs, _, err := client.Must().Apigateway.APIGatewaysApi.ApigatewaysGet(context.Background()).Execute()

	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(z apigateway.GatewayRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete gateway with name: %s, id: %s ", z.Properties.Name, z.Id),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().Apigateway.APIGatewaysApi.ApigatewaysDelete(c.Context, z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", z.Id, z.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
