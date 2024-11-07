package gateway

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	vpn "github.com/ionos-cloud/sdk-go-vpn"
)

var (
	allCols = []string{"ID", "Name", "Description", "GatewayIP", "InterfaceIPv4", "InterfaceIPv6", "DatacenterId", "LanId", "ConnectionIPv4", "ConnectionIPv6", "InterfaceIP", "ListenPort", "Status"}
	// we can safely include both InterfaceIPv4 and InterfaceIPv6 cols because if the respective column empty, it won't be shown
	defaultCols = []string{"ID", "Name", "Description", "GatewayIP", "InterfaceIPv4", "InterfaceIPv6", "InterfaceIP", "ListenPort", "Status"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "gateway",
			Short:            "Manage Wireguard VPN Gateways",
			Aliases:          []string{"g", "gw"},
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Get())
	cmd.AddCommand(Update())

	return cmd
}

// GatewaysProperty returns a list of properties of all gateways matching the given filters
func GatewaysProperty[V any](f func(gateway vpn.WireguardGatewayRead) V, fs ...Filter) []V {
	recs, err := Gateways(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Gateways returns all distributions matching the given filters
func Gateways(fs ...Filter) (vpn.WireguardGatewayReadList, error) {
	if url := config.GetServerUrl(); url == constants.DefaultApiURL || url == "" {
		viper.Set(constants.ArgServerUrl, constants.DefaultVPNApiURL)
	}

	req := client.Must().VPNClient.WireguardGatewaysApi.WireguardgatewaysGet(context.Background())
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vpn.WireguardGatewayReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vpn.WireguardGatewayReadList{}, err
	}
	return ls, nil
}

type Filter func(request vpn.ApiWireguardgatewaysGetRequest) (vpn.ApiWireguardgatewaysGetRequest, error)
