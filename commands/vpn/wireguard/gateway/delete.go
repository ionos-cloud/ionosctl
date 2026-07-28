package gateway

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/completer"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "vpn",
		Resource:  "wireguard gateway",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a gateway",
		Example:   "ionosctl vpn wg gateway delete " + core.FlagsUsage(constants.FlagGatewayID),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagGatewayID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			if err := c.RequireExplicitLocation(); err != nil {
				return err
			}

			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagGatewayID))

			g, _, err := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysFindById(context.Background(), id).Execute()
			if err != nil {
				return fmt.Errorf("failed getting gateway by id %s: %w", id, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf(
				"Are you sure you want to delete gateway %s at %s",
				g.Properties.Name, g.Properties.GatewayIP),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysDelete(context.Background(), id).Execute()

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagGatewayID, constants.FlagIdShort, "", "The ID of the WireGuard Gateway",
		core.RequiredFlagOption(),
		core.WithCompletion(completer.GatewayIDs, constants.VPNApiRegionalURL, constants.VPNLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all gateways. Required or --%s", constants.FlagGatewayID))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	// Gather gateways from every location (unless --location pins one), tagging each with its
	// location and location-scoped client, then hand the flat list to core.DeleteAll for a
	// consistent preview / per-item confirm-skip / summary flow.
	type located struct {
		gateway vpn.WireguardGatewayRead
		loc     string
		api     *vpn.APIClient
	}
	var items []located
	if err := c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		vc := vpn.NewAPIClient(cfg)
		xs, _, err := vc.WireguardGatewaysApi.WireguardgatewaysGet(context.Background()).Execute()
		if err != nil {
			return fmt.Errorf("failed listing gateways: %w", err)
		}
		for _, g := range xs.GetItems() {
			items = append(items, located{gateway: g, loc: location, api: vc})
		}
		return nil
	}); err != nil {
		return err
	}

	return core.DeleteAll(c, core.DeleteAllOptions[located]{
		Resource: "gateway",
		List:     func() ([]located, error) { return items, nil },
		Summary: func(l located) string {
			return fmt.Sprintf("%s (id: %s, ip: %s, location: %s)", l.gateway.Properties.Name, l.gateway.Id, l.gateway.Properties.GatewayIP, l.loc)
		},
		ID: func(l located) string { return l.gateway.Id },
		Delete: func(l located) error {
			_, err := l.api.WireguardGatewaysApi.WireguardgatewaysDelete(context.Background(), l.gateway.Id).Execute()
			return err
		},
	})
}
