package gateway

import (
	"context"
	"fmt"
	"net"
	"os"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"

	// "github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec gateway",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a IPSec Gateway",
		Example:   "ionosctl vpn ipsec gateway create " + core.FlagsUsage(constants.FlagName, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagConnectionIP, constants.FlagGatewayIP, constants.FlagInterfaceIP, constants.FlagPrivateKey),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			baseReq := []string{
				constants.FlagName,
				constants.FlagDatacenterId,
				constants.FlagLanId,
				constants.FlagConnectionIP,
				constants.FlagGatewayIP,
				constants.FlagInterfaceIP,
			}
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				// either privateKey or privateKeyPath are required
				append(baseReq, constants.FlagPrivateKey),
				append(baseReq, constants.FlagPrivateKeyPath),
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := &vpn.IPSecGateway{}

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

			if fn := core.GetFlagName(c.NS, constants.FlagPrivateKeyPath); viper.IsSet(fn) {
				// read the file
				keyBytes, err := os.ReadFile(viper.GetString(fn))
				if err != nil {
					return fmt.Errorf("failed to read private key file: %w", err)
				}
				input.PrivateKey = pointer.From(string(keyBytes))
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

			if fn := core.GetFlagName(c.NS, constants.FlagGatewayIP); viper.IsSet(fn) {
				input.GatewayIP = pointer.From(viper.GetString(fn))
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

			if fn := core.GetFlagName(c.NS, constants.FlagInterfaceIP); viper.IsSet(fn) {
				ip := viper.GetString(fn)
				if isIPv4(ip) {
					input.InterfaceIPv4CIDR = pointer.From(ip)
				} else {
					input.InterfaceIPv6CIDR = pointer.From(ip)
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

			createdGateway, _, err := client.Must().VPNClient.IPSecGatewaysApi.
				IPSecgatewaysPost(context.Background()).
				IPSecGatewayCreate(vpn.IPSecGatewayCreate{Properties: input}).Execute()
			if err != nil {
				return err
			}

			table, err := resource2table.ConvertVPNIPSecGatewayToTable(createdGateway)
			if err != nil {
				return fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutputPreconverted(createdGateway, table,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the IPSec Gateway", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Gateway")
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
	cmd.AddStringFlag(constants.FlagInterfaceIP, "", "", "The IPv4 or IPv6 address (with CIDR mask) to be assigned to the IPSec interface", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInterfaceIP, completer.GetCidrCompletionFunc(cmd))
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to connect your VPN Gateway to", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		loc, _ := c.Flags().GetString(constants.FlagLocation)
		return cloudapiv6completer.DatacenterIdsFilterLocation(loc), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID to connect your VPN Gateway to", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagConnectionIP, "", "", "A LAN IPv4 or IPv6 address in CIDR notation that will be assigned to the VPN Gateway", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagConnectionIP, completer.GetCidrCompletionFunc(cmd))

	cmd.AddStringFlag(constants.FlagPrivateKey, "K", "", fmt.Sprintf("Specify the private key (required or --%s)", constants.FlagPrivateKeyPath))
	cmd.AddStringFlag(constants.FlagPrivateKeyPath, "k", "", fmt.Sprintf("Specify the private key from a file (required or --%s)", constants.FlagPrivateKey))

	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that IPSec Server will listen on")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
