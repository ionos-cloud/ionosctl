package peer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/gateway"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard peer",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List WireGuard Peers",
		Example:   "ionosctl vpn wireguard peer list " + core.FlagsUsage(constants.FlagGatewayID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := Peers(
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID)),
				func(req vpn.ApiWireguardgatewaysPeersGetRequest) (vpn.ApiWireguardgatewaysPeersGetRequest, error) {
					if fn := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(fn) {
						req = req.Offset(viper.GetInt32(fn))
					}
					if fn := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(fn) {
						req = req.Limit(viper.GetInt32(fn))
					}
					return req, nil
				},
			)
			if err != nil {
				return fmt.Errorf("failed listing peers: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.VPNWireguardPeer, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return gateway.GatewaysProperty(func(gateway vpn.WireguardGatewayRead) string {
			return *gateway.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
