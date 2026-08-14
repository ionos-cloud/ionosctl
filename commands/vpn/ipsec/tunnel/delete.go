package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "ipsec tunnel",
		Verb:      "delete",
		Aliases:   []string{"d", "del", "rm"},
		ShortDesc: "Remove a IPSec Tunnel",
		Example:   "ionosctl vpn ipsec tunnel delete " + core.FlagsUsage(constants.FlagGatewayID, constants.FlagTunnelID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsSetsAndLocation(
				[]string{constants.FlagGatewayID, constants.FlagTunnelID},
				[]string{constants.FlagGatewayID, constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))
			p, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(), gatewayId, id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting tunnel by id %s: %w", id, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete tunnel %s"+
				" (host: '%s')", p.Properties.Name, p.Properties.RemoteHost),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsDelete(context.Background(), gatewayId, id).Execute()

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGatewayID, "", "", "The ID of the IPSec Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)
	cmd.AddStringFlag(constants.FlagTunnelID, constants.FlagIdShort, "", "The ID of the IPSec Tunnel you want to delete",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			gatewayID := viper.GetString(core.GetFlagName(cmd.NS, constants.FlagGatewayID))
			return completer.TunnelIDs(gatewayID)
		}, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all tunnels. Required or --%s", constants.FlagTunnelID))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))

	return core.DeleteAll(c, core.DeleteAllOptions[vpn.IPSecTunnelRead]{
		Resource: "tunnel",
		List: func() ([]vpn.IPSecTunnelRead, error) {
			xs, _, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsGet(context.Background(), gatewayId).Execute()
			if err != nil {
				return nil, err
			}
			return xs.GetItems(), nil
		},
		Summary: func(p vpn.IPSecTunnelRead) string {
			return fmt.Sprintf("%s (id: %s, host: %s)", p.Properties.Name, p.Id, p.Properties.RemoteHost)
		},
		ID: func(p vpn.IPSecTunnelRead) string { return p.Id },
		Delete: func(p vpn.IPSecTunnelRead) error {
			_, err := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsDelete(context.Background(), gatewayId, p.Id).Execute()
			return err
		},
	})
}
