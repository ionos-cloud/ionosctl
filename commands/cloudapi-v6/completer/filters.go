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

func DataCentersFiltersUsage() []string {
	return getPropertiesName(ionoscloud.DatacenterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func DataCentersPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.DatacenterProperties{})
}

func DataCentersMetadataFilters() []string {
	return getPropertiesName(ionoscloud.DatacenterElementMetadata{})
}

func BackupUnitsFilters() []string {
	return getPropertiesName(ionoscloud.BackupUnitProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func BackupUnitsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.BackupUnitProperties{})
}

func ServersFilters() []string {
	return getPropertiesName(ionoscloud.ServerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ServersPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.ServerProperties{})
}

func ImagesFilters() []string {
	return getPropertiesName(ionoscloud.ImageProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ImagesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.ImageProperties{})
}

func VolumesFilters() []string {
	return getPropertiesName(ionoscloud.VolumeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func VolumesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.VolumeProperties{})
}

func SnapshotsFilters() []string {
	return getPropertiesName(ionoscloud.SnapshotProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func SnapshotsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.SnapshotProperties{})
}

func IpBlocksFilters() []string {
	return getPropertiesName(ionoscloud.IpBlockProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func IpBlocksPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.IpBlockProperties{})
}

func LocationsFilters() []string {
	return getPropertiesName(ionoscloud.LocationProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LocationsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.LocationProperties{})
}

func LANsFilters() []string {
	return getPropertiesName(ionoscloud.LanProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LANsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.LanProperties{})
}

func NICsFilters() []string {
	return getPropertiesName(ionoscloud.NicProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NICsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.NicProperties{})
}

func FirewallRulesFilters() []string {
	return getPropertiesName(ionoscloud.FirewallruleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FirewallRulesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.FirewallruleProperties{})
}

func LoadBalancersFilters() []string {
	return getPropertiesName(ionoscloud.LoadbalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LoadBalancersPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.LoadbalancerProperties{})
}

func RequestsFilters() []string {
	return getPropertiesName(ionoscloud.RequestProperties{}, ionoscloud.RequestMetadata{})
}

func RequestsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.RequestProperties{})
}

func RequestsMetadataFilters() []string {
	return getPropertiesName(ionoscloud.RequestMetadata{})
}

func UsersFilters() []string {
	return getPropertiesName(ionoscloud.UserProperties{}, ionoscloud.UserMetadata{})
}

func UserPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.UserProperties{})
}

func UserMetadataFilters() []string {
	return getPropertiesName(ionoscloud.UserMetadata{})
}

func K8sClustersFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesClusterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sClustersPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesClusterProperties{})
}

func K8sNodePoolsFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodePoolProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodePoolsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodePoolProperties{})
}

func K8sNodesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodeProperties{})
}

func FlowLogsFilters() []string {
	return getPropertiesName(ionoscloud.FlowLogProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FlowLogsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.FlowLogProperties{})
}

func GroupsFilters() []string {
	return getPropertiesName(ionoscloud.GroupProperties{})
}

func NATGatewaysFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NATGatewaysPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayProperties{})
}

func NATGatewayRulesFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NATGatewayRulesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayRuleProperties{})
}

func NlbsFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NlbsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerProperties{})
}

func NlbForwardingRulesFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NlbForwardinRulesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerForwardingRuleProperties{})
}

func PccsFilters() []string {
	return getPropertiesName(ionoscloud.PrivateCrossConnectProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func PccsPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.PrivateCrossConnectProperties{})
}

func TemplatesFilters() []string {
	return getPropertiesName(ionoscloud.TemplateProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func TemplatesPropertiesFilters() []string {
	return getPropertiesName(ionoscloud.TemplateProperties{})
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
					} else {
						// Add slices of strings only, not nested properties as they are currently not supported
						if arg.Field(i).Type.Elem().String() == reflect.SliceOf(reflect.TypeOf("")).String() {
							properties = append(properties, makeFirstLowerCase(arg.Field(i).Name))
						}
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
