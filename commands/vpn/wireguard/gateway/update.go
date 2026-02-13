package gateway

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	dbaascompleter "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	// "github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/cobra"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "update",
		Aliases:   []string{"u", "put", "patch"},
		ShortDesc: "Update a WireGuard Gateway",
		LongDesc:  "Update a WireGuard Gateway. Note: The private key MUST be provided again (or changed) on updates.",
		Example:   "ionosctl vpn wireguard gateway update " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagName, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagConnectionIP, constants.FlagGatewayIP, constants.FlagInterfaceIP),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGatewayID, constants.FlagPrivateKey},
				[]string{constants.FlagGatewayID, constants.FlagPrivateKeyPath},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			id, err := c.Command.Command.Flags().GetString(constants.FlagGatewayID)
			if err != nil {
				return err
			}

			g, _, err := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysFindById(context.Background(), id).Execute()
			if err != nil {
				return err
			}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				g.Properties.Name = name
			}
			if c.Command.Command.Flags().Changed(constants.FlagDescription) {
				desc, err := c.Command.Command.Flags().GetString(constants.FlagDescription)
				if err != nil {
					return err
				}
				g.Properties.Description = pointer.From(desc)
			}
			if c.Command.Command.Flags().Changed(constants.FlagIp) {
				ip, err := c.Command.Command.Flags().GetString(constants.FlagIp)
				if err != nil {
					return err
				}
				g.Properties.GatewayIP = ip
			}

			if c.Command.Command.Flags().Changed(constants.FlagPrivateKey) {
				key, err := c.Command.Command.Flags().GetString(constants.FlagPrivateKey)
				if err != nil {
					return err
				}
				g.Properties.PrivateKey = key
			}

			if c.Command.Command.Flags().Changed(constants.FlagPrivateKeyPath) {
				// read the file
				keyPath, err := c.Command.Command.Flags().GetString(constants.FlagPrivateKeyPath)
				if err != nil {
					return err
				}
				keyBytes, err := os.ReadFile(keyPath)
				if err != nil {
					return fmt.Errorf("failed to read private key file: %w", err)
				}
				g.Properties.PrivateKey = string(keyBytes)
			}

			if c.Command.Command.Flags().Changed(constants.FlagPort) {
				port, err := c.Command.Command.Flags().GetInt32(constants.FlagPort)
				if err != nil {
					return err
				}
				g.Properties.ListenPort = pointer.From(port)
			}

			if c.Command.Command.Flags().Changed(constants.FlagDatacenterId) {
				if g.Properties.Connections == nil {
					g.Properties.Connections = make([]vpn.Connection, 1)
				}
				dcID, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
				if err != nil {
					return err
				}
				g.Properties.Connections[0].DatacenterId = dcID
			}
			if c.Command.Command.Flags().Changed(constants.FlagLanId) {
				if g.Properties.Connections == nil {
					g.Properties.Connections = make([]vpn.Connection, 1)
				}
				lanID, err := c.Command.Command.Flags().GetString(constants.FlagLanId)
				if err != nil {
					return err
				}
				g.Properties.Connections[0].LanId = lanID
			}

			if c.Command.Command.Flags().Changed(constants.FlagGatewayIP) {
				gatewayIP, err := c.Command.Command.Flags().GetString(constants.FlagGatewayIP)
				if err != nil {
					return err
				}
				g.Properties.GatewayIP = gatewayIP
			}

			// Note: VPN Gateway handles IPv4 and IPv6 addresses separately for both InterfaceIP and Connections.IP
			// We will use the same flag for both ipv4 and ipv6 for both of them, work out what type (v4 or v6) it is,
			// and pass it to the API as the correct field (ipv4 or ipv6)
			isIPv4 := func(ip string) bool {
				if ipAddr, _, err := net.ParseCIDR(ip); err == nil {
					return ipAddr.To4() != nil
				}
				return net.ParseIP(ip).To4() != nil
			}

			if c.Command.Command.Flags().Changed(constants.FlagInterfaceIP) {
				ip, err := c.Command.Command.Flags().GetString(constants.FlagInterfaceIP)
				if err != nil {
					return err
				}
				if isIPv4(ip) {
					g.Properties.InterfaceIPv4CIDR = pointer.From(ip)
				} else {
					g.Properties.InterfaceIPv6CIDR = pointer.From(ip)
				}
			}

			if c.Command.Command.Flags().Changed(constants.FlagConnectionIP) {
				if g.Properties.Connections == nil {
					g.Properties.Connections = make([]vpn.Connection, 1)
				}
				ip, err := c.Command.Command.Flags().GetString(constants.FlagConnectionIP)
				if err != nil {
					return err
				}
				if isIPv4(ip) {
					g.Properties.Connections[0].Ipv4CIDR = ip
				} else {
					g.Properties.Connections[0].Ipv6CIDR = pointer.From(ip)
				}
			}

			createdGateway, _, err := client.Must().VPNClient.WireguardGatewaysApi.
				WireguardgatewaysPut(context.Background(), id).
				WireguardGatewayEnsure(vpn.WireguardGatewayEnsure{Id: id, Properties: g.Properties}).Execute()
			if err != nil {
				return fmt.Errorf("failed updating gateway: %w", err)
			}

			table, err := resource2table.ConvertVPNWireguardGatewayToTable(createdGateway)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(createdGateway, table,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the WireGuard Gateway", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the WireGuard Gateway")
	cmd.AddStringFlag(constants.FlagGatewayIP, "", "", "The IP of an IPBlock in the same location as the provided datacenter", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayIP, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		dc, _, _ := client.Must().CloudClient.DataCentersApi.
			DatacentersFindById(context.Background(), viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))).
			Execute()

		ipblocks, _, err := client.Must().CloudClient.IPBlocksApi.
			IpblocksGet(context.Background()).
			Filter("location", *dc.Properties.Location).
			Execute()
		if err != nil || ipblocks.Items == nil || len(*ipblocks.Items) == 0 {
			return nil, cobra.ShellCompDirectiveError
		}
		var ips []string
		for _, ipblock := range *ipblocks.Items {
			if ipblock.Properties.Ips != nil {
				ips = append(ips, *ipblock.Properties.Ips...)
			}
		}
		return ips, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagInterfaceIP, "", "", "The IPv4 or IPv6 address (with CIDR mask) to be assigned to the WireGuard interface", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInterfaceIP, dbaascompleter.GetCidrCompletionFunc(cmd))
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to connect your VPN Gateway to", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID to connect your VPN Gateway to", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagConnectionIP, "", "", "A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagConnectionIP, dbaascompleter.GetCidrCompletionFunc(cmd))

	cmd.AddStringFlag(constants.FlagPrivateKey, "K", "", fmt.Sprintf("Specify the private key (required or --%s)", constants.FlagPrivateKeyPath))
	cmd.AddStringFlag(constants.FlagPrivateKeyPath, "k", "", fmt.Sprintf("Specify the private key from a file (required or --%s)", constants.FlagPrivateKey))

	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that WireGuard Server will listen on")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
