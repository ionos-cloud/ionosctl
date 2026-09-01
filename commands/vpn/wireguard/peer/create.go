package peer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a WireGuard Peer",
		LongDesc: `Add a remote device (peer) to a WireGuard gateway.

Supply the peer's --public-key so the gateway trusts it, and --ips: the source subnets the peer may send through the tunnel (WireGuard routes by key + allowed IPs). Use "a.b.c.d/32" for a single host, or "0.0.0.0/0","::/0" to allow everything.

--host/--port tell the gateway where to reach the peer for outbound connections; for peers behind NAT that only dial in, point --host at any reachable address (WireGuard learns the real endpoint from incoming traffic).

There is a per-gateway limit on peers; see product documentation.`,
		Example: `ionosctl vpn wireguard peer create --gateway-id GATEWAY_ID --name my-laptop --public-key PUBLIC_KEY --ips 10.7.222.0/24 --host vpn.example.com
ionosctl vpn wireguard peer create --gateway-id GATEWAY_ID --name allow-all --public-key PUBLIC_KEY --ips 0.0.0.0/0,::/0 --host 203.0.113.5 --port 51820`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(
				constants.FlagGatewayID, constants.FlagName, constants.FlagIps, constants.FlagPublicKey, constants.FlagHost,
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := vpn.WireguardPeer{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				input.Description = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagIps); viper.IsSet(fn) {
				input.AllowedIPs = viper.GetStringSlice(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPublicKey); viper.IsSet(fn) {
				input.PublicKey = viper.GetString(fn)
			}

			input.Endpoint = &vpn.WireguardEndpoint{}
			if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
				input.Endpoint.Host = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPort); viper.IsSet(fn) {
				input.Endpoint.Port = pointer.From(viper.GetInt32(fn))
			}

			peer, _, err := client.Must().VPNClient.WireguardPeersApi.
				WireguardgatewaysPeersPost(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))).
				WireguardPeerCreate(vpn.WireguardPeerCreate{Properties: input}).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(peer)
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the WireGuard Peer", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the WireGuard Peer")
	cmd.AddStringSliceFlag(constants.FlagIps, "", []string{}, "Comma-separated CIDRs the peer is allowed to send through the tunnel (its allowed source IPs). Use \"a.b.c.d/32\" for a single host, or \"0.0.0.0/0\",\"::/0\" for all addresses", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIps, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"::/0"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagPublicKey, "", "", "The peer's WireGuard public key; the gateway trusts the device holding the matching private key", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagHost, "", "", "Hostname or IPv4 the gateway uses to reach this peer (for peers behind NAT, any reachable address; the real endpoint is learned from inbound traffic)", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagPort, "", 51820, "UDP port the gateway uses to reach this peer (default 51820)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
