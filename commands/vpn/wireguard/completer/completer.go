package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
)

// -- GATEWAYS

// GatewaysProperty returns a list of properties of all gateways matching the given filters
func GatewaysProperty[V any](f func(gateway vpn.WireguardGatewayRead) V, fs ...GatewayFilter) []V {
	recs, err := Gateways(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

// Gateways returns all distributions matching the given filters
func Gateways(fs ...GatewayFilter) (vpn.WireguardGatewayReadList, error) {
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

type GatewayFilter func(request vpn.ApiWireguardgatewaysGetRequest) (vpn.ApiWireguardgatewaysGetRequest, error)

// -- PEERS

// PeersProperty returns a list of properties of all peers matching the given filters
func PeersProperty[V any](gatewayID string, f func(peer vpn.WireguardPeerRead) V, fs ...PeerFilter) []V {
	recs, err := Peers(gatewayID, fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

// Peers returns all distributions matching the given filters
func Peers(gatewayID string, fs ...PeerFilter) (vpn.WireguardPeerReadList, error) {
	req := client.Must().VPNClient.WireguardPeersApi.WireguardgatewaysPeersGet(context.Background(), gatewayID)
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vpn.WireguardPeerReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vpn.WireguardPeerReadList{}, err
	}
	return ls, nil
}

type PeerFilter func(request vpn.ApiWireguardgatewaysPeersGetRequest) (vpn.ApiWireguardgatewaysPeersGetRequest, error)

func GatewayIDs() []string {
	return GatewaysProperty(func(gateway vpn.WireguardGatewayRead) string {
		return gateway.Id + "\t" + gateway.Properties.Name + "(" + gateway.Properties.GatewayIP + ")"
	})
}

func PeerIDs(gatewayID string) []string {
	return PeersProperty(gatewayID, func(p vpn.WireguardPeerRead) string {
		return p.Id + "\t" + p.Properties.Name + "(" + p.Properties.Endpoint.Host + ")"
	})
}
