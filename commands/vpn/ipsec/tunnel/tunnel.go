package tunnel

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
)

/*
An IPSec tunnel is the connection from an IPSec gateway to one remote site. Unlike WireGuard, IPSec negotiates the connection in two phases:
  - IKE (phase 1) authenticates the two peers (PSK or RSA) and builds a secure control channel.
  - ESP (phase 2) is the child channel that actually encrypts the traffic.
Both phases have a Diffie-Hellman group, encryption algorithm, integrity algorithm and a lifetime (rekey interval, seconds). cloudNetworkCIDRs (local IONOS LAN side) and peerNetworkCIDRs (remote side) decide which subnets are allowed to cross the tunnel.
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
			Use:   "tunnel",
			Short: "Manage IPSec VPN Tunnels",
			Long: `Manage IPSec VPN Tunnels.

A tunnel is the connection from an IPSec gateway (--gateway-id) to ONE remote site. It carries everything the two ends must agree on:
  - --host: the remote peer's public address (IPv4 or FQDN).
  - authentication: --auth-method PSK with a --psk-key (or RSA).
  - IKE (phase 1) crypto: authenticates the peers and builds the control channel.
  - ESP (phase 2) crypto: the child channel that encrypts your traffic.
    Each phase takes a Diffie-Hellman group, encryption + integrity algorithm and a lifetime (rekey interval, seconds).
  - --cloud-network-cidrs (local IONOS LAN side) and --peer-network-cidrs (remote side): which subnets may cross.

Both ends must use MATCHING crypto and mirrored CIDRs, or the tunnel will not come up.`,
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
