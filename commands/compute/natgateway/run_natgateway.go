package natgateway

import (
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
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunNatGatewayListAll(c)
	}
	natgateways, resp, err := c.CloudApiV6Services.NatGateways().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(natgateways.NatGateways)
}

func RunNatGatewayGet(c *core.CommandConfig) error {
	c.Verbose("NAT Gateway with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)))

	ng, resp, err := c.CloudApiV6Services.NatGateways().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
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
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	ng, resp, err := c.CloudApiV6Services.NatGateways().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
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
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
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
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
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

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
		input.SetPublicIps(publicIps)

		c.Verbose("Property PublicIps set: %v", publicIps)
	}

	return &resources.NatGatewayProperties{
		NatGatewayProperties: input,
	}
}

func DeleteAllNatgateways(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	c.Verbose(constants.DatacenterId, dcId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.NatGateway]{
		Resource: "NAT Gateway",
		List: func() ([]ionoscloud.NatGateway, error) {
			natGateways, _, err := c.CloudApiV6Services.NatGateways().List(dcId)
			if err != nil {
				return nil, err
			}
			items, ok := natGateways.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of NAT Gateway")
			}
			return *items, nil
		},
		Summary: func(ng ionoscloud.NatGateway) string {
			summary := ""
			if props, ok := ng.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil {
					summary += *name
				}
				if ips, ok := props.GetPublicIpsOk(); ok && ips != nil && len(*ips) > 0 {
					summary += fmt.Sprintf(" (public IPs: %v)", *ips)
				}
			}
			if id, ok := ng.GetIdOk(); ok && id != nil {
				summary += fmt.Sprintf(" (id: %s)", *id)
			}
			return summary
		},
		ID: func(ng ionoscloud.NatGateway) string {
			if id := ng.GetId(); id != nil {
				return *id
			}
			return ""
		},
		Delete: func(ng ionoscloud.NatGateway) error {
			resp, err := c.CloudApiV6Services.NatGateways().Delete(dcId, *ng.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
