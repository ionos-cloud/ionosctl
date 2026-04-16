package tunnel

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

/*
A IPSec Tunnel is any device (client, server, or another gateway) that participates in a IPSec VPN. Tunnels are identified by public/private key pairs.
IPSec does not need complex negotiation (like IPsec IKE phases). Once two Tunnels know each other's public keys and IP addresses, they can connect instantly.
IPSec is stateless: no persistent state is stored between connections, and packets are exchanged only when needed.
There is no session or tunnel establishment process like in IPsec. Instead, IPSec Tunnels exchange packets as needed without keeping an active session.
*/

var allCols = []table.Column{
	{Name: "ID", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Description", JSONPath: "properties.description", Default: true},
	{Name: "RemoteHost", JSONPath: "properties.remoteHost", Default: true},
	{Name: "AuthMethod", JSONPath: "properties.auth.method", Default: true},
	{Name: "PSKKey", JSONPath: "properties.auth.psk.key", Default: true},
	{Name: "IKEDiffieHellmanGroup", JSONPath: "properties.ike.diffieHellmanGroup"},
	{Name: "IKEEncryptionAlgorithm", JSONPath: "properties.ike.encryptionAlgorithm"},
	{Name: "IKEIntegrityAlgorithm", JSONPath: "properties.ike.integrityAlgorithm"},
	{Name: "IKELifetime", JSONPath: "properties.ike.lifetime"},
	{Name: "ESPDiffieHellmanGroup", JSONPath: "properties.esp.diffieHellmanGroup"},
	{Name: "ESPEncryptionAlgorithm", JSONPath: "properties.esp.encryptionAlgorithm"},
	{Name: "ESPIntegrityAlgorithm", JSONPath: "properties.esp.integrityAlgorithm"},
	{Name: "ESPLifetime", JSONPath: "properties.esp.lifetime"},
	{Name: "CloudNetworkCIDRs", JSONPath: "properties.cloudNetworkCIDRs"},
	{Name: "PeerNetworkCIDRs", JSONPath: "properties.peerNetworkCIDRs"},
	{Name: "Status", JSONPath: "metadata.status", Default: true},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "tunnel",
			Short:            "Manage IPSec VPN Tunnels",
			Aliases:          []string{"p"},
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(Create())
	cmd.AddCommand(List())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Update())

	return cmd
}
