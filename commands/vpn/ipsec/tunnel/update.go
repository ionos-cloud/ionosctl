package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/gateway"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "update",
		Aliases:   []string{"u", "patch", "put"},
		ShortDesc: "Update a IPSec Tunnel",
		Example:   "", // TODO: Probably best if I don't forget this
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagTunnelID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))
			p, _, err := client.Must().VPNClient.IPSecTunnelsApi.IPSecgatewaysTunnelsFindById(context.Background(), gatewayId, id).Execute()

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				p.Properties.Name = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
				p.Properties.Description = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagIps); viper.IsSet(fn) {
				p.Properties.AllowedIPs = pointer.From(viper.GetStringSlice(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPublicKey); viper.IsSet(fn) {
				p.Properties.PublicKey = pointer.From(viper.GetString(fn))
			}

			p.Properties.Endpoint = &vpn.IPSecEndpoint{}
			if fn := core.GetFlagName(c.NS, constants.FlagHost); viper.IsSet(fn) {
				p.Properties.Endpoint.Host = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPort); viper.IsSet(fn) {
				p.Properties.Endpoint.Port = pointer.From(viper.GetInt32(fn))
			}

			tunnel, _, err := client.Must().VPNClient.IPSecTunnelsApi.
				IPSecgatewaysTunnelsPut(context.Background(), gatewayId, id).
				IPSecTunnelEnsure(vpn.IPSecTunnelEnsure{Id: &id, Properties: p.Properties}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.VPNIPSecTunnel, tunnel, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return gateway.GatewaysProperty(func(gateway vpn.IPSecGatewayRead) string {
			return *gateway.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel you want to delete", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagTunnelID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return TunnelsProperty(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagTunnelID)), func(p vpn.IPSecTunnelRead) string {
			return *p.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, "", "", "Name of the IPSec Tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "Description of the IPSec Tunnel")
	cmd.AddStringSliceFlag(constants.FlagIps, "", []string{}, "Comma separated subnets of CIDRs that are allowed to connect to the IPSec Gateway. Specify \"a.b.c.d/32\" for an individual IP address. Specify \"0.0.0.0/0\" or \"::/0\" for all addresses", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagIps, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"::/0"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagPublicKey, "", "", "Public key of the connecting tunnel", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagHost, "", "", "Hostname or IPV4 address that the IPSec Server will connect to", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagPort, "", 51820, "Port that the IPSec Server will connect to")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
