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
			Use:   "vpn",
			Short: "Manage VPN Gateways",
			Long: `Manage VPN Gateways: encrypted site-to-site tunnels between one of your IONOS CLOUD LANs and a remote network (an on-prem firewall, another cloud, a laptop).

A gateway lives in a datacenter and attaches to a LAN; it takes a public IP (from an IPBlock) that the remote side dials, and a private IP on that LAN so it can route cloud traffic into the tunnel. Two protocols are offered:

  wireguard  modern, key-based, connectionless. A gateway plus one 'peer' per remote device.
  ipsec      standards-based (IKE/ESP). A gateway plus one 'tunnel' per remote site, with negotiable crypto.

Pick WireGuard for simple, fast, key-pair setups; pick IPSec when you must interoperate with existing IPSec hardware/policies. Gateways are regional (--location) and limited to 8 per region.

Docs: https://docs.ionos.com/cloud/network-services/vpn-gateway`,
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(wireguard.Root())
	cmd.AddCommand(ipsec.Root())

	return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.VPN}, constants.VPNApiRegionalURL, constants.VPNLocations)
}
