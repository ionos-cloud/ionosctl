/*
This is used for supporting completion in the CLI.
Option: --filters
*/
package completer

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func DataCentersFilters() []string {
	return getPropertiesName(ionoscloud.DatacenterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func DataCentersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.DatacenterProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func BackupUnitsFilters() []string {
	return getPropertiesName(ionoscloud.BackupUnitProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func BackupUnitsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.BackupUnitProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func ServersFilters() []string {
	return getPropertiesName(ionoscloud.ServerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ServersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.ServerProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func ImagesFilters() []string {
	return getPropertiesName(ionoscloud.ImageProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ImagesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.ImageProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func VolumesFilters() []string {
	return getPropertiesName(ionoscloud.VolumeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func VolumesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.VolumeProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func SnapshotsFilters() []string {
	return getPropertiesName(ionoscloud.SnapshotProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func SnapshotsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.SnapshotProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func IpBlocksFilters() []string {
	return getPropertiesName(ionoscloud.IpBlockProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func IpBlocksFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.IpBlockProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func LabelsFilters() []string {
	return getPropertiesName(ionoscloud.LabelProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LabelsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.LabelProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func LocationsFilters() []string {
	return getPropertiesName(ionoscloud.LocationProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LocationsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.LocationProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func LANsFilters() []string {
	return getPropertiesName(ionoscloud.LanProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LANsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.LanProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func NICsFilters() []string {
	return getPropertiesName(ionoscloud.NicProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NICsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.NicProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func FirewallRulesFilters() []string {
	return getPropertiesName(ionoscloud.FirewallruleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FirewallRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.FirewallruleProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func ApplicationLoadBalancersFilters() []string {
	return getPropertiesName(ionoscloud.ApplicationLoadBalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func ApplicationLoadBalancersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.ApplicationLoadBalancerProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func LoadBalancersFilters() []string {
	return getPropertiesName(ionoscloud.LoadbalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func LoadbalancersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.LoadbalancerProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func RequestsFilters() []string {
	return getPropertiesName(ionoscloud.RequestProperties{}, ionoscloud.RequestMetadata{})
}

func RequestsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.RequestProperties{}), getPropertiesName(ionoscloud.RequestMetadata{}))
}

func UsersFilters() []string {
	return getPropertiesName(ionoscloud.UserProperties{}, ionoscloud.UserMetadata{})
}

func UsersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.UserProperties{}), getPropertiesName(ionoscloud.UserMetadata{}))
}

func K8sClustersFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesClusterProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sClustersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.KubernetesClusterProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func K8sNodePoolsFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodePoolProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodePoolsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.KubernetesNodePoolProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func K8sNodesFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesNodeProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func K8sNodesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.KubernetesNodeProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func FlowLogsFilters() []string {
	return getPropertiesName(ionoscloud.FlowLogProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func FlowLogsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.FlowLogProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func GroupsFilters() []string {
	return getPropertiesName(ionoscloud.GroupProperties{})
}

func GroupsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.GroupProperties{}), []string{})
}

func NATGatewaysFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NATGatewaysFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.NatGatewayProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func NATGatewayRulesFilters() []string {
	return getPropertiesName(ionoscloud.NatGatewayRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NATGatewayRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.NatGatewayRuleProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func NlbsFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NlbsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.NetworkLoadBalancerProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func NlbRulesFilters() []string {
	return getPropertiesName(ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func NlbRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.NetworkLoadBalancerForwardingRuleProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func PccsFilters() []string {
	return getPropertiesName(ionoscloud.PrivateCrossConnectProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func PccsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.PrivateCrossConnectProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
}

func TemplatesFilters() []string {
	return getPropertiesName(ionoscloud.TemplateProperties{}, ionoscloud.DatacenterElementMetadata{})
}

func TemplatesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(ionoscloud.TemplateProperties{}), getPropertiesName(ionoscloud.DatacenterElementMetadata{}))
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

func getFilterUsage(propertiesFilters []string, metadataFilters []string) string {
	usage := "Available Filters:\n"
	if len(propertiesFilters) > 0 {
		usage = fmt.Sprintf("%s* filter by property: %s", usage, propertiesFilters)
	}
	if len(metadataFilters) > 0 {
		usage = fmt.Sprintf("%s\n* filter by metadata: %s", usage, metadataFilters)
	}
	return usage
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
