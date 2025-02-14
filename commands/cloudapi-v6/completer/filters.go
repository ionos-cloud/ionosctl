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

	compute "github.com/ionos-cloud/sdk-go/v6"
)

func DataCentersFilters() []string {
	return getPropertiesName(compute.DatacenterProperties{}, compute.DatacenterElementMetadata{})
}

func DataCentersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.DatacenterProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func BackupUnitsFilters() []string {
	return getPropertiesName(compute.BackupUnitProperties{}, compute.DatacenterElementMetadata{})
}

func BackupUnitsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.BackupUnitProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func ServersFilters() []string {
	return getPropertiesName(compute.ServerProperties{}, compute.DatacenterElementMetadata{})
}

func ServersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.ServerProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func ImagesFilters() []string {
	return getPropertiesName(compute.ImageProperties{}, compute.DatacenterElementMetadata{})
}

func ImagesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.ImageProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func VolumesFilters() []string {
	return getPropertiesName(compute.VolumeProperties{}, compute.DatacenterElementMetadata{})
}

func VolumesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.VolumeProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func SnapshotsFilters() []string {
	return getPropertiesName(compute.SnapshotProperties{}, compute.DatacenterElementMetadata{})
}

func SnapshotsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.SnapshotProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func IpBlocksFilters() []string {
	return getPropertiesName(compute.IpBlockProperties{}, compute.DatacenterElementMetadata{})
}

func IpBlocksFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.IpBlockProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func LabelsFilters() []string {
	return getPropertiesName(compute.LabelProperties{}, compute.DatacenterElementMetadata{})
}

func LabelsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.LabelProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func LocationsFilters() []string {
	return getPropertiesName(compute.LocationProperties{}, compute.DatacenterElementMetadata{})
}

func LocationsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.LocationProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func LANsFilters() []string {
	return getPropertiesName(compute.LanProperties{}, compute.DatacenterElementMetadata{})
}

func LANsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.LanProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func NICsFilters() []string {
	return getPropertiesName(compute.NicProperties{}, compute.DatacenterElementMetadata{})
}

func NICsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.NicProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func FirewallRulesFilters() []string {
	return getPropertiesName(compute.FirewallruleProperties{}, compute.DatacenterElementMetadata{})
}

func FirewallRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.FirewallruleProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func ApplicationLoadBalancersFilters() []string {
	return getPropertiesName(compute.ApplicationLoadBalancerProperties{}, compute.DatacenterElementMetadata{})
}

func ApplicationLoadBalancersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.ApplicationLoadBalancerProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func LoadBalancersFilters() []string {
	return getPropertiesName(compute.LoadbalancerProperties{}, compute.DatacenterElementMetadata{})
}

func LoadbalancersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.LoadbalancerProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func RequestsFilters() []string {
	return getPropertiesName(compute.RequestProperties{}, compute.RequestMetadata{}, compute.RequestStatusMetadata{})
}

func RequestsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.RequestProperties{}), getPropertiesName(compute.RequestMetadata{}, compute.RequestStatusMetadata{}))
}

func UsersFilters() []string {
	return getPropertiesName(compute.UserProperties{}, compute.UserMetadata{})
}

func UsersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.UserProperties{}), getPropertiesName(compute.UserMetadata{}))
}

func K8sClustersFilters() []string {
	return getPropertiesName(compute.KubernetesClusterProperties{}, compute.DatacenterElementMetadata{})
}

func K8sClustersFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.KubernetesClusterProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func K8sNodePoolsFilters() []string {
	return getPropertiesName(compute.KubernetesNodePoolProperties{}, compute.DatacenterElementMetadata{})
}

func K8sNodePoolsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.KubernetesNodePoolProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func K8sNodesFilters() []string {
	return getPropertiesName(compute.KubernetesNodeProperties{}, compute.DatacenterElementMetadata{})
}

func K8sNodesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.KubernetesNodeProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func FlowLogsFilters() []string {
	return getPropertiesName(compute.FlowLogProperties{}, compute.DatacenterElementMetadata{})
}

func FlowLogsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.FlowLogProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func GroupsFilters() []string {
	return getPropertiesName(compute.GroupProperties{})
}

func GroupsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.GroupProperties{}), []string{})
}

func NATGatewaysFilters() []string {
	return getPropertiesName(compute.NatGatewayProperties{}, compute.DatacenterElementMetadata{})
}

func NATGatewaysFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.NatGatewayProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func NATGatewayRulesFilters() []string {
	return getPropertiesName(compute.NatGatewayRuleProperties{}, compute.DatacenterElementMetadata{})
}

func NATGatewayRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.NatGatewayRuleProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func NlbsFilters() []string {
	return getPropertiesName(compute.NetworkLoadBalancerProperties{}, compute.DatacenterElementMetadata{})
}

func NlbsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.NetworkLoadBalancerProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func NlbRulesFilters() []string {
	return getPropertiesName(compute.NetworkLoadBalancerForwardingRuleProperties{}, compute.DatacenterElementMetadata{})
}

func NlbRulesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.NetworkLoadBalancerForwardingRuleProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func PccsFilters() []string {
	return getPropertiesName(compute.PrivateCrossConnectProperties{}, compute.DatacenterElementMetadata{})
}

func PccsFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.PrivateCrossConnectProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
}

func TemplatesFilters() []string {
	return getPropertiesName(compute.TemplateProperties{}, compute.DatacenterElementMetadata{})
}

func TemplatesFiltersUsage() string {
	return getFilterUsage(getPropertiesName(compute.TemplateProperties{}), getPropertiesName(compute.DatacenterElementMetadata{}))
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
