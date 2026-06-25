package gpu

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func PreRunServerGpusList(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)

}

func RunServerGpusList(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	serverId := c.Flags().String(cloudapiv6.ArgServerId)

	gpus, _, err := client.Must().CloudClient.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(c.Context, dcId, serverId).Execute()
	if err != nil {
		return fmt.Errorf("failed to list Gpus from Server %s: %w", serverId, err)
	}

	return c.Printer(allGpuCols).Prefix("items").Print(gpus)
}

func PreRunDcServerGpuIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgGpuId)
}

func RunServerGpuGet(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	serverId := c.Flags().String(cloudapiv6.ArgServerId)
	gpuId := c.Flags().String(cloudapiv6.ArgGpuId)

	gpu, _, err := client.Must().CloudClient.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsFindById(c.Context, dcId, serverId, gpuId).Execute()
	if err != nil {
		return fmt.Errorf("failed to get GPU from Server %s: %w", serverId, err)
	}

	return c.Printer(allGpuCols).Print(gpu)
}
