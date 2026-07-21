package natgateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunNATGatewayList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunDcIdsNatGatewayIps(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIps)
}

func PreRunDcNatGatewayIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId)
}

func PreRunNatGatewayDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)
	var allNatGateways []ionoscloud.NatGateways
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter ID")
		}

		natGateways, resp, err := c.CloudApiV6Services.NatGateways().List(*id)
		if err != nil {
			return err
		}

		allNatGateways = append(allNatGateways, natGateways.NatGateways)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allCols).Prefix("*.items").Print(allNatGateways)
}

func RunNatGatewayList(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunNatGatewayListAll(c)
	}
	natgateways, resp, err := c.CloudApiV6Services.NatGateways().List(c.Flags().String(cloudapiv6.ArgDataCenterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(natgateways.NatGateways)
}

func RunNatGatewayGet(c *core.CommandConfig) error {
	c.Verbose("NAT Gateway with id: %v is getting...", c.Flags().String(cloudapiv6.ArgNatGatewayId))

	ng, resp, err := c.CloudApiV6Services.NatGateways().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNatGatewayId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGateway)
}

func RunNatGatewayCreate(c *core.CommandConfig) error {
	proper := getNewNatGatewayInfo(c)

	if !proper.HasName() {
		proper.SetName(c.Flags().String(cloudapiv6.ArgName))
	}

	ng, resp, err := c.CloudApiV6Services.NatGateways().Create(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		resources.NatGateway{
			NatGateway: ionoscloud.NatGateway{
				Properties: &proper.NatGatewayProperties,
			},
		},
	)

	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGateway)
}

func RunNatGatewayUpdate(c *core.CommandConfig) error {
	input := getNewNatGatewayInfo(c)

	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNatGatewayId),
		*input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGateway)
}

func RunNatGatewayDelete(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	natGatewayId := c.Flags().String(cloudapiv6.ArgNatGatewayId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllNatgateways(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nat gateway", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starring deleting NAT Gateway with id: %v...", natGatewayId)

	resp, err := c.CloudApiV6Services.NatGateways().Delete(dcId, natGatewayId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("NAT Gateway successfully deleted")
	return nil
}

func getNewNatGatewayInfo(c *core.CommandConfig) *resources.NatGatewayProperties {
	input := ionoscloud.NatGatewayProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgIps) {
		publicIps := c.Flags().StringSlice(cloudapiv6.ArgIps)
		input.SetPublicIps(publicIps)

		c.Verbose("Property PublicIps set: %v", publicIps)
	}

	return &resources.NatGatewayProperties{
		NatGatewayProperties: input,
	}
}

func DeleteAllNatgateways(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting NatGateways...")

	natGateways, resp, err := c.CloudApiV6Services.NatGateways().List(dcId)
	if err != nil {
		return err
	}

	natGatewayItems, ok := natGateways.GetItemsOk()
	if !ok || natGatewayItems == nil {
		return fmt.Errorf("could not get items of NAT Gateway")
	}

	if len(*natGatewayItems) <= 0 {
		return fmt.Errorf("no NAT Gateways found")
	}

	var multiErr error
	for _, natGateway := range *natGatewayItems {
		name := natGateway.GetProperties().Name
		id := natGateway.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the NAT Gateway with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NatGateways().Delete(dcId, *id)
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
