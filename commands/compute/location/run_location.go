package location

import (
	"errors"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgLocationId)
}

func RunLocationList(c *core.CommandConfig) error {

	locations, resp, err := c.CloudApiV6Services.Locations().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLocationCols).Prefix("items").Print(locations)
}

func RunLocationGet(c *core.CommandConfig) error {
	locId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId))
	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}

	c.Verbose("Location with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId)))

	loc, resp, err := c.CloudApiV6Services.Locations().GetByRegionAndLocationId(ids[0], ids[1])
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLocationCols).Print(loc)
}
