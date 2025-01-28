package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
)

// -- GATEWAYS

// GatewaysProperty returns a list of properties of all gateways matching the given filters
func GatewaysProperty[V any](f func(gateway vpn.IPSecGatewayRead) V, fs ...GatewayFilter) []V {
	recs, err := Gateways(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

// Gateways returns all distributions matching the given filters
func Gateways(fs ...GatewayFilter) (vpn.IPSecGatewayReadList, error) {
	req := client.Must().VPNClient.IPSecGatewaysApi.IpsecgatewaysGet(context.Background())
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vpn.IPSecGatewayReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vpn.IPSecGatewayReadList{}, err
	}
	return ls, nil
}

type GatewayFilter func(request vpn.ApiIpsecgatewaysGetRequest) (vpn.ApiIpsecgatewaysGetRequest, error)

// TunnelsProperty returns a list of properties of all tunnels matching the given filters
func TunnelsProperty[V any](gatewayID string, f func(tunnel vpn.IPSecTunnelRead) V, fs ...TunnelFilter) []V {
	recs, err := Tunnels(gatewayID, fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

// Tunnels returns all distributions matching the given filters
func Tunnels(gatewayID string, fs ...TunnelFilter) (vpn.IPSecTunnelReadList, error) {
	req := client.Must().VPNClient.IPSecTunnelsApi.IpsecgatewaysTunnelsGet(context.Background(), gatewayID)
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

type TunnelFilter func(request vpn.ApiIpsecgatewaysTunnelsGetRequest) (vpn.ApiIpsecgatewaysTunnelsGetRequest, error)

func GatewayIDs() []string {
	return GatewaysProperty(func(gateway vpn.IPSecGatewayRead) string {
		return gateway.Id + "\t" + gateway.Properties.Name + "(" + gateway.Properties.GatewayIP + ")"
	})
}

func TunnelIDs(gatewayID string) []string {
	return TunnelsProperty(gatewayID, func(p vpn.IPSecTunnelRead) string {
		return p.Id + "\t" + p.Properties.Name + "(" + p.Properties.RemoteHost + ")"
	})
}
