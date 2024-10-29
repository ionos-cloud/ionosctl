package gateway

import (
	"context"
	"fmt"
	"net"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a WireGuard Gateway",
		Example:   "", // TODO: Probably best if I don't forget this
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagName,
				constants.FlagIp,
				constants.FlagInterfaceIP,
				constants.FlagDatacenterId,
				constants.FlagLanId,
				constants.FlagConnectionIP,
				constants.FlagPrivateKey,
				constants.FlagPort,
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := &vpn.WireguardGateway{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				input.Description = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagIp); viper.IsSet(fn) {
				input.GatewayIP = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPrivateKey); viper.IsSet(fn) {
				input.PrivateKey = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPort); viper.IsSet(fn) {
				input.ListenPort = pointer.From(viper.GetInt32(fn))
			}

			input.Connections = pointer.From(make([]vpn.Connection, 1))
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				(*input.Connections)[0].DatacenterId = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				(*input.Connections)[0].LanId = pointer.From(viper.GetString(fn))
			}

			// Note: VPN Gateway handles IPv4 and IPv6 addresses separately for both InterfaceIP and Connections.IP
			// We will use the same flag for both ipv4 and ipv6 for both of them, work out what type (v4 or v6) it is,
			// and pass it to the API as the correct field (ipv4 or ipv6)

			isIPv4 := func(ip string) bool {
				return net.ParseIP(ip).To4() != nil
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInterfaceIP); viper.IsSet(fn) {
				ip := viper.GetString(fn)
				if isIPv4(ip) {
					input.InterfaceIPv4CIDR = pointer.From(ip)
				} else {
					input.InterfaceIPv4CIDR = pointer.From(ip)
				}
			}

			if fn := core.GetFlagName(c.NS, constants.FlagConnectionIP); viper.IsSet(fn) {
				ip := viper.GetString(fn)
				if isIPv4(ip) {
					(*input.Connections)[0].Ipv4CIDR = pointer.From(ip)
				} else {
					(*input.Connections)[0].Ipv6CIDR = pointer.From(ip)
				}
			}

			createdGateway, _, err := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysPost(context.Background()).
				WireguardGatewayCreate(vpn.WireguardGatewayCreate{Properties: input}).Execute()
			if err != nil {
				return err
			}

			json, err := createdGateway.MarshalJSON()
			if err != nil {
				return err
			}

			fmt.Println(string(json))

			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the WireGuard Gateway", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the WireGuard Gateway")
	cmd.AddStringFlag(constants.FlagGatewayIP, "", "", "Public IP address to be assigned to the gateway. Note: This must be an IP address in the same datacenter as the connections", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagInterfaceIP, "", "", "The IPv4 or IPv6 address (with CIDR mask) to be assigned to the WireGuard interface", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInterfaceIP, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ipblocks, _, err := client.Must().CloudClient.IPBlocksApi.IpblocksGet(context.Background()).Execute()
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
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagConnectionIP, completer.GetCidrCompletionFunc(cmd))

	cmd.AddStringFlag(constants.FlagPrivateKey, "k", "", "The private key to be used by the WireGuard Gateway", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that WireGuard Server will listen on")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
