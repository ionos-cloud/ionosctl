package tunnel

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
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
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGatewayID, constants.FlagTunnelID},
				[]string{constants.FlagGatewayID, constants.ArgAll},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			vpnClient := vpn.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))

			gatewayId := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagTunnelID))
			p, _, err := vpnClient.IPSecTunnelsApi.IpsecgatewaysTunnelsFindById(context.Background(), gatewayId, id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting tunnel by id %s: %w", id, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete tunnel %s"+
				" (host: '%s')", p.Properties.Name, p.Properties.RemoteHost),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = vpnClient.IPSecTunnelsApi.IpsecgatewaysTunnelsDelete(context.Background(), gatewayId, id).Execute()

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
	c.Verbose("Deleting all tunnels from gateway %s!", gatewayId)

	vpnClient := vpn.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))

	xs, _, err := vpnClient.IPSecTunnelsApi.IpsecgatewaysTunnelsGet(context.Background(), gatewayId).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(p vpn.IPSecTunnelRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf(
			"Are you sure you want to delete tunnel %s at %s",
			p.Properties.Name, p.Properties.RemoteHost),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := vpnClient.IPSecGatewaysApi.IpsecgatewaysDelete(context.Background(), p.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", p.Id, p.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
