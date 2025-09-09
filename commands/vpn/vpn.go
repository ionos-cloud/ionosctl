package vpn

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/ipsec"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "vpn",
			Short:            "VPN Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(wireguard.Root())
	cmd.AddCommand(ipsec.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.VPN}, constants.VPNApiRegionalURL, constants.VPNLocations)
}
