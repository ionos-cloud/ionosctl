package tunnel

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

/*
A IPSec Tunnel is any device (client, server, or another gateway) that participates in a IPSec VPN. Tunnels are identified by public/private key pairs.
IPSec does not need complex negotiation (like IPsec IKE phases). Once two Tunnels know each other’s public keys and IP addresses, they can connect instantly.
IPSec is stateless: no persistent state is stored between connections, and packets are exchanged only when needed.
There is no session or tunnel establishment process like in IPsec. Instead, IPSec Tunnels exchange packets as needed without keeping an active session.
*/

/*
	var VPNIPSecTunnel = map[string]string{
		"ID":                     "id",
		"Name":                   "properties.name",
		"Description":            "properties.description",
		"RemoteHost":             "properties.remoteHost",
		"AuthMethod":             "properties.auth.method",
		"PSKKey":                 "properties.auth.psk.key",
		"IKEDiffieHellmanGroup":  "properties.ike.diffieHellmanGroup",
		"IKEEncryptionAlgorithm": "properties.ike.encryptionAlgorithm",
		"IKEIntegrityAlgorithm":  "properties.ike.integrityAlgorithm",
		"IKELifetime":            "properties.ike.lifetime",
		"ESPDiffieHellmanGroup":  "properties.esp.diffieHellmanGroup",
		"ESPEncryptionAlgorithm": "properties.esp.encryptionAlgorithm",
		"ESPIntegrityAlgorithm":  "properties.esp.integrityAlgorithm",
		"ESPLifetime":            "properties.esp.lifetime",
		"CloudNetworkCIDRs":      "properties.cloudNetworkCIDRs",
		"PeerNetworkCIDRs":       "properties.peerNetworkCIDRs",
	}
*/
var (
	allCols = []string{"ID", "Name", "Description", "RemoteHost", "AuthMethod", "PSKKey",
		"IKEDiffieHellmanGroup", "IKEEncryptionAlgorithm", "IKEIntegrityAlgorithm", "IKELifetime",
		"ESPDiffieHellmanGroup", "ESPEncryptionAlgorithm", "ESPIntegrityAlgorithm", "ESPLifetime",
		"CloudNetworkCIDRs", "PeerNetworkCIDRs", "Status", "StatusMessage"}
	defaultCols = []string{"ID", "Name", "Description", "RemoteHost", "AuthMethod", "PSKKey", "Status"}
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

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Update())

	return cmd
}
