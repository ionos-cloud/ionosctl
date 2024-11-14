package completer

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/gateway"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn/wireguard/peer"
	vpn "github.com/ionos-cloud/sdk-go-vpn"
)

func GatewayIDs() []string {
	return gateway.GatewaysProperty(func(gateway vpn.WireguardGatewayRead) string {
		return *gateway.Id + "\t" + *gateway.Properties.Name + "(" + *gateway.Properties.GatewayIP + ")"
	})
}

func PeerIDs(gatewayID string) []string {
	return peer.PeersProperty(gatewayID, func(p vpn.WireguardPeerRead) string {
		return *p.Id + "\t" + *p.Properties.Name + "(" + *p.Properties.Endpoint.Host + ")"
	})
}
