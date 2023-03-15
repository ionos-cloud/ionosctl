package completer

import (
	"context"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"io"

	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// DatacenterCPUFamilies returns the list of available CPU families in a given datacenter
// If the datacenter ID is not provided returns a set of possible default values.
func DatacenterCPUFamilies(ctx context.Context, outErr io.Writer, datacenterId string) []string {
	if datacenterId == "" {
		return []string{"AMD_OPTERON", "INTEL_XEON", "INTEL_SKYLAKE"}
	}
	client, err := client2.Get()
	clierror.CheckError(err, outErr)
	dcSvc := resources.NewDataCenterService(client, ctx)
	dc, _, err := dcSvc.Get(datacenterId, resources.QueryParams{})
	clierror.CheckError(err, outErr)
	if dc.Properties == nil || dc.Properties.CpuArchitecture == nil {
		return nil
	}
	result := make([]string, len(*dc.Properties.CpuArchitecture))
	for i, arch := range *dc.Properties.CpuArchitecture {
		result[i] = *arch.CpuFamily
	}
	return result
}
