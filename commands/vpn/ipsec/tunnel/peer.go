package tunnel

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*
A IPSec Tunnel is any device (client, server, or another gateway) that participates in a IPSec VPN. Tunnels are identified by public/private key pairs.
IPSec does not need complex negotiation (like IPsec IKE phases). Once two Tunnels know each other’s public keys and IP addresses, they can connect instantly.
IPSec is stateless: no persistent state is stored between connections, and packets are exchanged only when needed.
There is no session or tunnel establishment process like in IPsec. Instead, IPSec Tunnels exchange packets as needed without keeping an active session.
*/

var (
	allCols = []string{"ID", "Name", "Description", "Host", "Port", "WhitelistIPs", "PublicKey", "Status"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "tunnel",
			Short:            "Manage IPSec VPN Tunnels",
			Aliases:          []string{"p"},
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Update())

	return cmd
}

// TunnelsProperty returns a list of properties of all tunnels matching the given filters
func TunnelsProperty[V any](gatewayID string, f func(tunnel vpn.IPSecTunnelRead) V, fs ...Filter) []V {
	recs, err := Tunnels(gatewayID, fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Tunnels returns all distributions matching the given filters
func Tunnels(gatewayID string, fs ...Filter) (vpn.IPSecTunnelReadList, error) {
	if url := config.GetServerUrl(); url == constants.DefaultApiURL || url == "" {
		viper.Set(constants.ArgServerUrl, constants.DefaultVPNApiURL)
	}

	req := client.Must().VPNClient.IPSecTunnelsApi.IPSecgatewaysTunnelsGet(context.Background(), gatewayID)
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vpn.IPSecTunnelReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vpn.IPSecTunnelReadList{}, err
	}
	return ls, nil
}

type Filter func(request vpn.ApiIPSecgatewaysTunnelsGetRequest) (vpn.ApiIPSecgatewaysTunnelsGetRequest, error)
