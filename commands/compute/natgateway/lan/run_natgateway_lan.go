package lan

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunDcNatGatewayLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgLanId)
}

func PreRunDcNatGatewayLanRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayLanList(c *core.CommandConfig) error {
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

	return c.Printer(allCols).Prefix("properties.lans").Print(ng.NatGateway)
}

func RunNatGatewayLanAdd(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	natGatewayId := c.Flags().String(cloudapiv6.ArgNatGatewayId)

	ng, _, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	c.Verbose("Adding NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId)

	input := getNewNatGatewayLanInfo(c, ng)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("properties.lans").Print(ng.NatGateway)
}

func RunNatGatewayLanRemove(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := RemoveAllNatGatewayLans(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove nat gateway lan", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	natGatewayId := c.Flags().String(cloudapiv6.ArgNatGatewayId)

	ng, _, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	c.Verbose("Removing NatGateway with id %v to Datacenter with id: %v", natGatewayId, dcId)

	input := removeNatGatewayLanInfo(c, ng)
	ng, resp, err := c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("NAT Gateway Lan successfully deleted")
	return nil
}

func RemoveAllNatGatewayLans(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	natGatewayId := c.Flags().String(cloudapiv6.ArgNatGatewayId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("NatGateway ID: %v", natGatewayId)
	c.Verbose("Getting NatGateway...")

	natGateway, resp, err := c.CloudApiV6Services.NatGateways().Get(dcId, natGatewayId)
	if err != nil {
		return err
	}

	natGatewayProperties, ok := natGateway.GetPropertiesOk()
	if !ok || natGatewayProperties == nil {
		return fmt.Errorf("could not get NAT Gateway properties")
	}

	lansOk, ok := natGatewayProperties.GetLansOk()
	if !ok || lansOk == nil {
		return fmt.Errorf("could not get items of NAT Gateway Lans")
	}

	if len(*lansOk) <= 0 {
		return fmt.Errorf("no NAT Gateway Lans found")
	}

	c.Msg("NAT Gateway Lans to be removed:")
	for _, lanItem := range *lansOk {
		if id, ok := lanItem.GetIdOk(); ok && id != nil {
			c.Msg("NAT Gateway Lan Id: %v", string(*id))
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the NAT Gateways Lans", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Removing all the NAT Gateways Lans...")

	proper := make([]ionoscloud.NatGatewayLanProperties, 0)
	if natGateway != nil {
		if properties, ok := natGateway.GetPropertiesOk(); ok && properties != nil {
			natGatewaysProps := &resources.NatGatewayProperties{
				NatGatewayProperties: ionoscloud.NatGatewayProperties{
					Lans: &proper,
				},
			}

			natGateway, resp, err = c.CloudApiV6Services.NatGateways().Update(dcId, natGatewayId, *natGatewaysProps)
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			if err != nil {
				return err
			}

		}
	}

	c.Msg("NAT Gateway Lans successfully deleted")
	return nil
}

func getNewNatGatewayLanInfo(c *core.CommandConfig, oldNg *resources.NatGateway) *resources.NatGatewayProperties {
	var proper []ionoscloud.NatGatewayLanProperties

	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				proper = *lans
			}
		}
	}

	input := ionoscloud.NatGatewayLanProperties{}
	if c.Flags().Changed(cloudapiv6.ArgLanId) {
		lanId := c.Flags().Int32(cloudapiv6.ArgLanId)
		input.SetId(lanId)

		c.Verbose("Property Id set: %v", lanId)
	}

	if c.Flags().Changed(cloudapiv6.ArgIps) {
		gatewayIps := c.Flags().StringSlice(cloudapiv6.ArgIps)
		input.SetGatewayIps(gatewayIps)

		c.Verbose("Property GatewayIps set: %v", gatewayIps)
	}

	proper = append(proper, input)

	return &resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

func removeNatGatewayLanInfo(c *core.CommandConfig, oldNg *resources.NatGateway) *resources.NatGatewayProperties {
	proper := make([]ionoscloud.NatGatewayLanProperties, 0)

	if oldNg != nil {
		if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					if id, ok := lanItem.GetIdOk(); ok && id != nil {
						if *id != c.Flags().Int32(cloudapiv6.ArgLanId) {
							proper = append(proper, lanItem)
						}
					}
				}
			}
		}
	}

	return &resources.NatGatewayProperties{
		NatGatewayProperties: ionoscloud.NatGatewayProperties{
			Lans: &proper,
		},
	}
}

func PreRunDcNatGatewayIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId)
}
