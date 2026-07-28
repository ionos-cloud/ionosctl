package rule

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunNATGatewayRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId); err != nil {
		return err
	}
	return nil
}

func PreRunNatGatewayRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgIp, cloudapiv6.ArgSourceSubnet)
}

func PreRunDcNatGatewayRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgRuleId)
}

func PreRunDcNatGatewayRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgAll},
	)
}

func RunNatGatewayRuleList(c *core.CommandConfig) error {
	natgatewayRules, resp, err := c.CloudApiV6Services.NatGateways().ListRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(natgatewayRules.NatGatewayRules)
}

func RunNatGatewayRuleGet(c *core.CommandConfig) error {
	c.Verbose("NatGatewayRule with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	ng, resp, err := c.CloudApiV6Services.NatGateways().GetRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGatewayRule)
}

func RunNatGatewayRuleCreate(c *core.CommandConfig) error {
	proper := getNewNatGatewayRuleInfo(c)

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if !proper.HasProtocol() {
		proper.SetProtocol(ionoscloud.NatGatewayRuleProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))
		c.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}

	ng, resp, err := c.CloudApiV6Services.NatGateways().CreateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		resources.NatGatewayRule{
			NatGatewayRule: ionoscloud.NatGatewayRule{
				Properties: &proper.NatGatewayRuleProperties,
			},
		},
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGatewayRule)
}

func RunNatGatewayRuleUpdate(c *core.CommandConfig) error {
	input := getNewNatGatewayRuleInfo(c)
	ng, resp, err := c.CloudApiV6Services.NatGateways().UpdateRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		*input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NatGatewayRule)
}

func RunNatGatewayRuleDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNatgatewayRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nat gateway rule", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting NatGatewayRule with id: %v...", ruleId)

	resp, err := c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, ruleId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("NAT Gateway Rule successfully deleted")
	return nil
}

func getNewNatGatewayRuleInfo(c *core.CommandConfig) *resources.NatGatewayRuleProperties {
	input := ionoscloud.NatGatewayRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIp)) {
		publicIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
		input.SetPublicIp(publicIp)

		c.Verbose("Property PublicIp set: %v", publicIp)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		protocol := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		input.SetProtocol(ionoscloud.NatGatewayRuleProtocol(protocol))

		c.Verbose("Property Protocol set: %v", protocol)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceSubnet)) {
		sourceSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceSubnet))
		input.SetSourceSubnet(sourceSubnet)

		c.Verbose("Property SourceSubnet set: %v", sourceSubnet)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetSubnet)) {
		targetSubnet := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetSubnet))
		input.SetTargetSubnet(targetSubnet)

		c.Verbose("Property Name set: %v", targetSubnet)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart)) &&
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd)) {
		inputPortRange := ionoscloud.TargetPortRange{}

		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart))
		portRangeStop := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd))

		inputPortRange.SetStart(portRangeStart)
		inputPortRange.SetEnd(portRangeStop)
		input.SetTargetPortRange(inputPortRange)

		c.Verbose("Property TargetPortRange set with start: %v and stop: %v", portRangeStart, portRangeStop)
	}

	return &resources.NatGatewayRuleProperties{
		NatGatewayRuleProperties: input,
	}
}

func DeleteAllNatgatewayRules(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natGatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("NatGateway ID: %v", natGatewayId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.NatGatewayRule]{
		Resource: "NAT Gateway Rule",
		List: func() ([]ionoscloud.NatGatewayRule, error) {
			natGatewayRules, _, err := c.CloudApiV6Services.NatGateways().ListRules(dcId, natGatewayId)
			if err != nil {
				return nil, err
			}
			items, ok := natGatewayRules.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of NAT Gateway Rules")
			}
			return *items, nil
		},
		Summary: func(rule ionoscloud.NatGatewayRule) string {
			summary := ""
			if props, ok := rule.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil {
					summary += *name
				}
				if ip, ok := props.GetPublicIpOk(); ok && ip != nil && *ip != "" {
					summary += fmt.Sprintf(" (public IP: %s)", *ip)
				}
			}
			if id, ok := rule.GetIdOk(); ok && id != nil {
				summary += fmt.Sprintf(" (id: %s)", *id)
			}
			return summary
		},
		ID: func(rule ionoscloud.NatGatewayRule) string {
			if id := rule.GetId(); id != nil {
				return *id
			}
			return ""
		},
		Delete: func(rule ionoscloud.NatGatewayRule) error {
			resp, err := c.CloudApiV6Services.NatGateways().DeleteRule(dcId, natGatewayId, *rule.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
