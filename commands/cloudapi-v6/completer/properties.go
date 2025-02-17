package completer

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// DatacenterCPUFamilies returns the list of available CPU families in a given datacenter
// If the datacenter ID is not provided returns a set of possible default values.
func DatacenterCPUFamilies(ctx context.Context, datacenterId string) []string {
	if datacenterId == "" {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}
	}
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	dcSvc := resources.NewDataCenterService(client, ctx)
	dc, _, err := dcSvc.Get(datacenterId, resources.QueryParams{})
	if err != nil {
		return nil
	}
	if dc.Properties == nil || dc.Properties.CpuArchitecture == nil {
		return nil
	}
	result := make([]string, len(*dc.Properties.CpuArchitecture))
	for i, arch := range dc.Properties.CpuArchitecture {
		result[i] = *arch.CpuFamily
	}
	return result
}
