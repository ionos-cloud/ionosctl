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

func K8sClusterFilters() []string {
	return getPropertiesName(ionoscloud.KubernetesClusterProperties{}, ionoscloud.DatacenterElementMetadata{})
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
