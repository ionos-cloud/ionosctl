package peer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Find a peer by ID",
		Example:   "ionosctl vpn wg peer get " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagPeerID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagPeerID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagPeerID))

			p, _, err := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersFindById(context.Background(), gatewayId, id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting peer by id %s: %w", id, err)
			}

			return c.Printer(allCols).Print(p)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)
	cmd.AddStringFlag(constants.FlagPeerID, constants.FlagIdShort, "", "The ID of the WireGuard Peer",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.PeerIDs(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID)))
		}, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
