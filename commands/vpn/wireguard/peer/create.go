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
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a WireGuard Peer",
		LongDesc:  "Create WireGuard Peers. There is a limit to the total number of peers. Please refer to product documentation",
		Example:   "ionosctl vpn wireguard peer create " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagName, constants.FlagIps, constants.FlagPublicKey, constants.FlagHost),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagGatewayID, constants.FlagName, constants.FlagIps, constants.FlagPublicKey, constants.FlagHost,
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := vpn.WireguardPeer{}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				input.Name = name
			}

			if c.Command.Command.Flags().Changed(constants.FlagDescription) {
				desc, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
				if err != nil {
					return err
				}
				input.Description = pointer.From(desc)
			}

			if c.Command.Command.Flags().Changed(constants.FlagIps) {
				ips, err := c.Command.Command.Flags().GetStringSlice(constants.FlagIps)
				if err != nil {
					return err
				}
				input.AllowedIPs = ips
			}

			if c.Command.Command.Flags().Changed(constants.FlagPublicKey) {
				key, err := c.Command.Command.Flags().GetString(constants.FlagPublicKey)
				if err != nil {
					return err
				}
				input.PublicKey = key
			}

			input.Endpoint = &vpn.WireguardEndpoint{}
			if c.Command.Command.Flags().Changed(constants.FlagHost) {
				host, err := c.Command.Command.Flags().GetString(constants.FlagHost)
				if err != nil {
					return err
				}
				input.Endpoint.Host = host
			}

			if c.Command.Command.Flags().Changed(constants.FlagPort) {
				port, err := c.Command.Command.Flags().GetInt32(constants.FlagPort)
				if err != nil {
					return err
				}
				input.Endpoint.Port = pointer.From(port)
			}

			gatewayID, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}
			peer, _, err := client.Must().VPNClient.WireguardPeersApi.
				WireguardgatewaysPeersPost(context.Background(), gatewayID).
				WireguardPeerCreate(vpn.WireguardPeerCreate{Properties: input}).Execute()
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

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
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
