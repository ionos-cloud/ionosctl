package client

import (
	"bytes"
	"reflect"
	"strings"
	"sync"

	cloudv6 "github.com/ionos-cloud/sdk-go/v6"
)

var (
	filterKeyMap  map[string]string // lowercase → correct camelCase
	filterKeyOnce sync.Once
)

// normalizeFilterKey maps any-cased filter key to the correct camelCase form
// expected by the IONOS Cloud API. If the key is unknown, it is returned as-is.
// E.g. "IMaGeTYpE" → "imageType", "imagetype" → "imageType", "imageType" → "imageType"
func normalizeFilterKey(key string) string {
	filterKeyOnce.Do(func() {
		filterKeyMap = buildFilterKeyMap()
	})

	if canonical, ok := filterKeyMap[strings.ToLower(key)]; ok {
		return canonical
	}
	return key
}

func buildFilterKeyMap() map[string]string {
	m := make(map[string]string)

	// All known CloudAPI v6 property and metadata struct types.
	// These must match the types used in commands/compute/completer/filters.go.
	types := []any{
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

	for _, t := range types {
		rt := reflect.TypeOf(t)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < rt.NumField(); i++ {
			camelCase := firstLower(rt.Field(i).Name)
			m[strings.ToLower(camelCase)] = camelCase
		}
	}

	return m
}

// firstLower converts "ImageType" → "imageType" (Go exported name → JSON camelCase).
func firstLower(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	return string(bytes.Join([][]byte{
		bytes.ToLower([]byte{bts[0]}), bts[1:],
	}, nil))
}
