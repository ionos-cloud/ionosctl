package peer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "update",
		Aliases:   []string{"u", "patch", "put"},
		ShortDesc: "Update a WireGuard Peer",
		Example:   "ionosctl vpn wireguard peer update " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagPeerID, constants.FlagName, constants.FlagDescription, constants.FlagIps, constants.FlagPublicKey, constants.FlagHost, constants.FlagPort),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagPeerID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			gatewayId, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}
			id, err := c.Command.Command.Flags().GetString(constants.FlagPeerID)
			if err != nil {
				return err
			}
			p, _, err := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersFindById(context.Background(), gatewayId, id).Execute()

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				p.Properties.Name = name
			}

			if c.Command.Command.Flags().Changed(constants.FlagDescription) {
				desc, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
				if err != nil {
					return err
				}
				p.Properties.Description = pointer.From(desc)
			}

			if c.Command.Command.Flags().Changed(constants.FlagIps) {
				ips, err := c.Command.Command.Flags().GetStringSlice(constants.FlagIps)
				if err != nil {
					return err
				}
				p.Properties.AllowedIPs = ips
			}

			if c.Command.Command.Flags().Changed(constants.FlagPublicKey) {
				key, err := c.Command.Command.Flags().GetString(constants.FlagPublicKey)
				if err != nil {
					return err
				}
				p.Properties.PublicKey = key
			}

			if c.Command.Command.Flags().Changed(constants.FlagHost) {
				if p.Properties.Endpoint == nil {
					p.Properties.Endpoint = &vpn.WireguardEndpoint{}
				}
				host, err := c.Command.Command.Flags().GetString(constants.FlagHost)
				if err != nil {
					return err
				}
				p.Properties.Endpoint.Host = host
			}

			if c.Command.Command.Flags().Changed(constants.FlagPort) {
				if p.Properties.Endpoint == nil {
					p.Properties.Endpoint = &vpn.WireguardEndpoint{}
				}
				port, err := c.Command.Command.Flags().GetInt32(constants.FlagPort)
				if err != nil {
					return err
				}
				p.Properties.Endpoint.Port = pointer.From(port)
			}

			peer, _, err := client.Must().VPNClient.WireguardPeersApi.
				WireguardgatewaysPeersPut(context.Background(), gatewayId, id).
				WireguardPeerEnsure(vpn.WireguardPeerEnsure{Id: id, Properties: p.Properties}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.VPNWireguardPeer, peer, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

			return nil
		},
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the WireGuard Peer", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the WireGuard Peer")
	cmd.AddStringSliceFlag(constants.FlagIps, "", []string{}, "Comma separated subnets of CIDRs that are allowed to connect to the WireGuard Gateway. Specify \"a.b.c.d/32\" for an individual IP address. Specify \"0.0.0.0/0\" or \"::/0\" for all addresses", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIps, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"::/0"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagPublicKey, "", "", "Public key of the connecting peer", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagHost, "", "", "Hostname or IPV4 address that the WireGuard Server will connect to", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that the WireGuard Server will connect to")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
