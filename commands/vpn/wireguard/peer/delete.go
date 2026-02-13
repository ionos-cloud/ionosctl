package peer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "delete",
		Aliases:   []string{"d", "del", "rm"},
		ShortDesc: "Remove a WireGuard Peer",
		Example:   "ionosctl vpn wireguard peer delete " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagPeerID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGatewayID, constants.FlagPeerID},
				[]string{constants.FlagGatewayID, constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			all, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
			if err != nil {
				return err
			}
			if all {
				return deleteAll(c)
			}

			gatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}
			id, err := c.Command.Command.Flags().GetString(constants.FlagPeerID)
			if err != nil {
				return err
			}
			p, _, err := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersFindById(context.Background(), gatewayId, id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting peer by id %s: %w", id, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf(
				"Are you sure you want to delete peer %s (host: '%s')",
				p.Properties.Name, p.Properties.Endpoint.Host),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersDelete(context.Background(), gatewayId, id).Execute()

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)
	cmd.AddStringFlag(constants.FlagPeerID, constants.FlagIdShort, "", "The ID of the WireGuard Peer you want to delete",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.PeerIDs(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID)))
		}, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all peers. Required or --%s", constants.FlagPeerID))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	gatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all peers from gateway %s!", gatewayId))

	xs, _, err := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersGet(context.Background(), gatewayId).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(p vpn.WireguardPeerRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf(
			"Are you sure you want to delete peer %s at %s", p.Properties.Name, p.Properties.Endpoint.Host),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysDelete(context.Background(), p.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", p.Id, p.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
