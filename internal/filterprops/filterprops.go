// Package filterprops provides shared CloudAPI v6 filter property helpers
// used by both the client filter-key normalizer and the CLI completer.
package filterprops

import (
	"bytes"
	"strings"

	cloudv6 "github.com/ionos-cloud/sdk-go/v6"
)

// FirstLower converts a Go exported name to JSON camelCase by lowering the
// first character: "ImageType" → "imageType".
func FirstLower(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	return string(bytes.Join([][]byte{
		bytes.ToLower([]byte{bts[0]}), bts[1:],
	}, nil))
}

// AllCloudAPIv6Types returns the canonical list of CloudAPI v6 SDK property
// and metadata struct types whose fields are valid filter keys.
func AllCloudAPIv6Types() []any {
	return []any{
		cloudv6.DatacenterProperties{},
		cloudv6.DatacenterElementMetadata{},
		cloudv6.ServerProperties{},
		cloudv6.ImageProperties{},
		cloudv6.VolumeProperties{},
		cloudv6.SnapshotProperties{},
		cloudv6.IpBlockProperties{},
		cloudv6.LabelProperties{},
		cloudv6.LocationProperties{},
		cloudv6.LanProperties{},
		cloudv6.NicProperties{},
		cloudv6.FirewallruleProperties{},
		cloudv6.ApplicationLoadBalancerProperties{},
		cloudv6.LoadbalancerProperties{},
		cloudv6.RequestProperties{},
		cloudv6.RequestMetadata{},
		cloudv6.RequestStatusMetadata{},
		cloudv6.UserProperties{},
		cloudv6.UserMetadata{},
		cloudv6.KubernetesClusterProperties{},
		cloudv6.KubernetesNodePoolProperties{},
		cloudv6.KubernetesNodeProperties{},
		cloudv6.FlowLogProperties{},
		cloudv6.GroupProperties{},
		cloudv6.NatGatewayProperties{},
		cloudv6.NatGatewayRuleProperties{},
		cloudv6.NetworkLoadBalancerProperties{},
		cloudv6.NetworkLoadBalancerForwardingRuleProperties{},
		cloudv6.PrivateCrossConnectProperties{},
		cloudv6.TemplateProperties{},
		cloudv6.BackupUnitProperties{},
	}
}
