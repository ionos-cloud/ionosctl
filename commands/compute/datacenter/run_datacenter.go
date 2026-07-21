package datacenter

import (
	"errors"
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
	c.Verbose("Getting Datacenter with ID: %v...", c.Flags().String(cloudapiv6.ArgDataCenterId))

	dc, resp, err := c.CloudApiV6Services.DataCenters().Get(c.Flags().String(cloudapiv6.ArgDataCenterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allDatacenterCols).Print(dc.Datacenter)
}

func RunDataCenterCreate(c *core.CommandConfig) error {
	name := c.Flags().String(cloudapiv6.ArgName)
	description := c.Flags().String(cloudapiv6.ArgDescription)
	loc := c.Flags().String(cloudapiv6.ArgLocation)

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

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)
		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgDescription) {
		description := c.Flags().String(cloudapiv6.ArgDescription)
		input.SetDescription(description)
		c.Verbose("Property Description set: %v", description)
	}

	dc, resp, err := c.CloudApiV6Services.DataCenters().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
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
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllDatacenters(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete data center", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

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
	c.Verbose("Getting Datacenters...")

	datacenters, resp, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	datacentersItems, ok := datacenters.GetItemsOk()
	if !ok || datacentersItems == nil {
		return fmt.Errorf("could not get items of Datacenters")
	}

	if len(*datacentersItems) <= 0 {
		return fmt.Errorf("no Datacenters found")
	}

	c.Msg("Datacenters to be deleted:")

	var multiErr error
	for _, dc := range *datacentersItems {
		id := dc.GetId()
		name := dc.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Datacenter with Id: %s , Name: %s", *id, *name), viper.IsSet(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.DataCenters().Delete(*id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func getDataCenters(datacenters resources.Datacenters) []resources.Datacenter {
	dc := make([]resources.Datacenter, 0)
	if items, ok := datacenters.GetItemsOk(); ok && items != nil {
		for _, datacenter := range *items {
			dc = append(dc, resources.Datacenter{Datacenter: datacenter})
		}
	}
	return dc
}

func GetIPv6CidrBlockFromDatacenter(dc ionoscloud2.Datacenter) (string, error) {
	if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
		if ipv6CidrBlock, ok := properties.GetIpv6CidrBlockOk(); ok && ipv6CidrBlock != nil {
			return *ipv6CidrBlock, nil
		}
	}

	return "", fmt.Errorf("could not get IPv6 Cidr Block from Datacenter")
}
