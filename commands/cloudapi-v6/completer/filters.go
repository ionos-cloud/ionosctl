/*
	This is used for supporting completion in the CLI.
	Option: --filters
*/
package completer

import (
	"bytes"
	"reflect"
	"strings"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func DataCentersFilters() []string {
	return getPropertiesName(ionoscloud.DatacenterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func BackupUnitFilters() []string {
	return getPropertiesName(ionoscloud.BackupUnitProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ServersFilters() []string {
	return getPropertiesName(ionoscloud.ServerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ImagesFilters() []string {
	return getPropertiesName(ionoscloud.ImageProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func VolumesFilters() []string {
	return getPropertiesName(ionoscloud.VolumeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func SnapshotsFilters() []string {
	return getPropertiesName(ionoscloud.SnapshotProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func IpBlocksFilters() []string {
	return getPropertiesName(ionoscloud.IpBlockProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LANsFilters() []string {
	return getPropertiesName(ionoscloud.LanProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NICsFilters() []string {
	return getPropertiesName(ionoscloud.NicProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FirewallRulesFilters() []string {
	return getPropertiesName(ionoscloud.FirewallruleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LoadBalancersFilters() []string {
	return getPropertiesName(ionoscloud.LoadbalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func RequestsFilters() []string {
	return getPropertiesName(ionoscloud.RequestProperties{}, ionoscloud.RequestMetadata{})
}

func UsersFilters() []string {
	return getPropertiesName(ionoscloud.UserProperties{}, ionoscloud.UserMetadata{})
}

func K8sClustersFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesClusterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodePoolsFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodePoolProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FlowLogsFilters() []string {
	return getPropertiesName(ionoscloud.FlowLogProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func GroupFilters() []string {
	return getPropertiesName(ionoscloud.GroupProperties{})
}

func NATGatewayFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NATGatewayRuleFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NetworkLoadBalancerFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NetworkLoadBalancerForwardingRulesFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func PccsFilters() []string {
	return getPropertiesName(ionoscloud.PrivateCrossConnectProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func TemplatesFilters() []string {
	return getPropertiesName(ionoscloud.TemplateProperties{}, ionoscloud.DatacenterElementMetadata{})
}


// getPropertiesName uses reflection to get properties' name from a struct.
// It helps in making the filters available to the user in autocompletion.
func getPropertiesName(params ...interface{}) []string {
	var argKind reflect.Kind
	properties := make([]string, 0)
	for _, param := range params {
		arg := reflect.TypeOf(param)
		if arg.Kind() == reflect.Ptr {
			argKind = arg.Elem().Kind()
		} else {
			argKind = arg.Kind()
		}
		if argKind == reflect.Struct {
			for i := 0; i < arg.NumField(); i++ {
				if arg.Field(i).Type.Kind() == reflect.Ptr {
					argKind = arg.Field(i).Type.Elem().Kind()
					if argKind != reflect.Slice {
						properties = append(properties, makeFirstLowerCase(arg.Field(i).Name))
					}
				}
			}
		}
	}
	return properties
}

func makeFirstLowerCase(s string) string {
	if len(s) < 2 {
		return strings.ToLower(s)
	}
	bts := []byte(s)
	return string(bytes.Join([][]byte{
		bytes.ToLower([]byte{bts[0]}), bts[1:],
	}, nil))
}
