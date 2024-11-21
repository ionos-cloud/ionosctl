package vpn

import (
	"fmt"
	"maps"
	"slices"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vpn",
			Short:            "VPN Operations",
			TraverseChildren: true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				location, _ := cmd.Flags().GetString(constants.FlagLocation)
				changeLocation(client.Must().VPNClient, location)
			},
		},
	}

	cmd.AddCommand(wireguard.Root())
	cmd.AddCommand(ipsec.Root())

	cmd.Command.PersistentFlags().String(constants.FlagLocation, "de/txl", fmt.Sprintf("The location your resources are hosted in. Possible values: %s", slices.Collect(maps.Keys(locationToURL))))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var locations []string
		for k := range locationToURL {
			locations = append(locations, k)
		}

		return locations, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

var locationToURL = map[string]string{
	"de/fra": "https://vpn.de-fra.ionos.com",
	"de/txl": "https://vpn.de-txl.ionos.com",
}

func changeLocation(client *vpn.APIClient, location string) {
	cfg := client.GetConfig()
	cfg.Servers = vpn.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
	return
}
