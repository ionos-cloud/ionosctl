package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a gateway",
		Example:   "ionosctl vpn wg g delete ...", // TODO: Probably best if I don't forget this
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagGatewayID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))

			g, _, err := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting gateway by id %s: %w", id, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete gateway %s (desc: '%s')", *g.Properties.Name, *g.Properties.Description),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysDelete(context.Background(), id).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return GatewaysProperty(func(gateway vpn.WireguardGatewayRead) string {
			return *gateway.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all zones. Required or -%s", constants.FlagZoneShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all gateways!"))
	xs, _, err := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(g vpn.WireguardGatewayRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete gateway %s at %s", *g.Properties.Name, *g.Properties.GatewayIP),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysDelete(context.Background(), *g.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", *g.Id, *g.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
