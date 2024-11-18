package tunnel

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Find a tunnel by ID",
		Example:   "ionosctl vpn ipsec tunnel get " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagTunnelID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagGatewayID, constants.FlagTunnelID)
		},
		CmdRun: func(c *core.CommandConfig) error {
			gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))

			p, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(), gatewayId, id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting tunnel by id %s: %w", id, err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.VPNIPSecTunnel, p,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the IPSec Gateway", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GatewayIDs(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel you want to delete", core.RequiredFlagOption())
	cmd.Command.RegisterFlagCompletionFunc(constants.FlagGatewayID, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TunnelIDs(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
