package datacenter

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud2 "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
}

func PreRunDataCenterDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunDataCenterList(c *core.CommandConfig) error {
	datacenters, resp, err := c.CloudApiV6Services.DataCenters().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allDatacenterCols).Prefix("items").Print(datacenters.Datacenters)
}

func RunDataCenterGet(c *core.CommandConfig) error {
	c.Verbose("Getting Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	dc, resp, err := c.CloudApiV6Services.DataCenters().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allDatacenterCols).Print(dc.Datacenter)
}

func RunDataCenterCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
	loc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))

	c.Verbose("Properties set for creating the datacenter: Name: %v, Description: %v, Location: %v", name, description, loc)

	dc, resp, err := c.CloudApiV6Services.DataCenters().Create(name, description, loc)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allDatacenterCols).Print(dc)
}

func RunDataCenterUpdate(c *core.CommandConfig) error {
	input := resources.DatacenterPropertiesPut{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
		description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))
		input.SetDescription(description)
		c.Verbose("Property Description set: %v", description)
	}

	dc, resp, err := c.CloudApiV6Services.DataCenters().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer(allDatacenterCols).Print(dc)
}

func RunDataCenterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllDatacenters(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete data center", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	c.Verbose("Starting deleting Datacenter with ID: %v...", dcId)

	resp, err := c.CloudApiV6Services.DataCenters().Delete(dcId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Datacenter successfully deleted")

	return nil

}

func DeleteAllDatacenters(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud2.Datacenter]{
		Resource: "datacenter",
		List: func() ([]ionoscloud2.Datacenter, error) {
			datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
			if err != nil {
				return nil, err
			}

			items, ok := datacenters.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Datacenters")
			}

			return *items, nil
		},
		Summary: func(dc ionoscloud2.Datacenter) string {
			var id, name, location, description string
			if dc.Id != nil {
				id = *dc.Id
			}
			if p := dc.Properties; p != nil {
				if p.Name != nil {
					name = *p.Name
				}
				if p.Location != nil {
					location = *p.Location
				}
				if p.Description != nil {
					description = *p.Description
				}
			}

			s := fmt.Sprintf("%s (id: %s, location: %s)", name, id, location)
			if description != "" {
				s = fmt.Sprintf("%s (id: %s, location: %s, desc: %s)", name, id, location, description)
			}
			return s
		},
		ID: func(dc ionoscloud2.Datacenter) string {
			if dc.Id != nil {
				return *dc.Id
			}
			return ""
		},
		Delete: func(dc ionoscloud2.Datacenter) error {
			resp, err := c.CloudApiV6Services.DataCenters().Delete(*dc.Id)
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

func GetIPv6CidrBlockFromDatacenter(dc ionoscloud2.Datacenter) (string, error) {
	if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
		if ipv6CidrBlock, ok := properties.GetIpv6CidrBlockOk(); ok && ipv6CidrBlock != nil {
			return *ipv6CidrBlock, nil
		}
	}

	return "", fmt.Errorf("could not get IPv6 Cidr Block from Datacenter")
}
