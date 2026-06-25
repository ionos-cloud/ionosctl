package location

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func RunLocationCpuList(c *core.CommandConfig) error {
	locId := c.Flags().String(cloudapiv6.ArgLocationId)

	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}

	loc, resp, err := c.CloudApiV6Services.Locations().GetByRegionAndLocationId(ids[0], ids[1])
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	properties, ok := loc.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting location properties")
	}

	cpus, ok := properties.GetCpuArchitectureOk()
	if !ok || cpus == nil {
		return fmt.Errorf("error getting cpu architectures")
	}

	return c.Printer(allCpuCols).Print(*cpus)
}
