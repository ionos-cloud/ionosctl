package firewallrule

import (
	"errors"
	"fmt"
	"net"
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

func PreRunFirewallRuleList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId); err != nil {
		return err
	}
	return nil
}

func PreRunDcServerNicIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId)
}

func PreRunFirewallDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgFirewallRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgAll},
	)
}

func PreRunDcServerNicIdsFRuleProtocol(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgProtocol)
}

func PreRunDcServerNicFRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgFirewallRuleId)
}

func RunFirewallRuleList(c *core.CommandConfig) error {
	firewallRules, resp, err := c.CloudApiV6Services.FirewallRules().List(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
		c.Flags().String(cloudapiv6.ArgNicId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFirewallRuleCols).Prefix("items").Print(firewallRules.FirewallRules)
}

func RunFirewallRuleGet(c *core.CommandConfig) error {
	c.Verbose("Getting Firewall Rule with id: %v", c.Flags().String(cloudapiv6.ArgFirewallRuleId))

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
		c.Flags().String(cloudapiv6.ArgNicId),
		c.Flags().String(cloudapiv6.ArgFirewallRuleId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFirewallRuleCols).Print(firewallRule.FirewallRule)
}

func RunFirewallRuleCreate(c *core.CommandConfig) error {
	if c.Flags().Changed(cloudapiv6.FlagIPVersion) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	properties := getFirewallRulePropertiesSet(c)

	if !properties.HasName() {
		properties.SetName(c.Flags().String(cloudapiv6.ArgName))
	}

	if !properties.HasType() {
		properties.SetType(c.Flags().String(cloudapiv6.ArgDirection))
	}

	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Create(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
		c.Flags().String(cloudapiv6.ArgNicId),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFirewallRuleCols).Print(firewallRule.FirewallRule)
}

func RunFirewallRuleUpdate(c *core.CommandConfig) error {
	if c.Flags().Changed(cloudapiv6.FlagIPVersion) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgServerId),
		c.Flags().String(cloudapiv6.ArgNicId),
		c.Flags().String(cloudapiv6.ArgFirewallRuleId),
		getFirewallRulePropertiesSet(c),
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFirewallRuleCols).Print(firewallRule.FirewallRule)
}

func RunFirewallRuleDelete(c *core.CommandConfig) error {
	datacenterId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	serverId := c.Flags().String(cloudapiv6.ArgServerId)
	nicId := c.Flags().String(cloudapiv6.ArgNicId)
	fruleId := c.Flags().String(cloudapiv6.ArgFirewallRuleId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllFirewallRules(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete Firewall Rule with id: %v...", fruleId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err := c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, fruleId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Firewall Rule successfully deleted")

	return nil

}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *core.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		properties.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgProtocol) {
		protocol := c.Flags().String(cloudapiv6.ArgProtocol)
		properties.SetProtocol(protocol)

		c.Verbose("Property Protocol set: %v", protocol)
	}

	if c.Flags().Changed(cloudapiv6.ArgSourceIp) {
		sourceIp := c.Flags().String(cloudapiv6.ArgSourceIp)
		properties.SetSourceIp(sourceIp)

		c.Verbose("Property SourceIp set: %v", sourceIp)
	}

	if c.Flags().Changed(cloudapiv6.ArgSourceMac) {
		sourceMac := c.Flags().String(cloudapiv6.ArgSourceMac)
		properties.SetSourceMac(sourceMac)

		c.Verbose("Property SourceMac set: %v", sourceMac)
	}

	if c.Flags().Changed(cloudapiv6.ArgDestinationIp) {
		targetIp := c.Flags().String(cloudapiv6.ArgDestinationIp)
		properties.SetTargetIp(targetIp)

		c.Verbose("Property TargetIp/DestinationIp set: %v", targetIp)
	}

	if c.Flags().Changed(cloudapiv6.ArgIcmpCode) {
		icmpCode := c.Flags().Int32(cloudapiv6.ArgIcmpCode)
		properties.SetIcmpCode(icmpCode)

		c.Verbose("Property IcmpCode set: %v", icmpCode)
	}

	if c.Flags().Changed(cloudapiv6.ArgIcmpType) {
		icmpType := c.Flags().Int32(cloudapiv6.ArgIcmpType)
		properties.SetIcmpType(icmpType)

		c.Verbose("Property IcmpType set: %v", icmpType)
	}

	if c.Flags().Changed(cloudapiv6.ArgPortRangeStart) {
		portRangeStart := c.Flags().Int32(cloudapiv6.ArgPortRangeStart)
		properties.SetPortRangeStart(portRangeStart)

		c.Verbose("Property PortRangeStart set: %v", portRangeStart)
	}

	if c.Flags().Changed(cloudapiv6.ArgPortRangeEnd) {
		portRangeEnd := c.Flags().Int32(cloudapiv6.ArgPortRangeEnd)
		properties.SetPortRangeEnd(portRangeEnd)

		c.Verbose("Property PortRangeEnd set: %v", portRangeEnd)
	}

	if c.Flags().Changed(cloudapiv6.ArgDirection) {
		firewallruleType := c.Flags().String(cloudapiv6.ArgDirection)
		properties.SetType(strings.ToUpper(firewallruleType))

		c.Verbose("Property Type/Direction set: %v", firewallruleType)
	}

	if c.Flags().Changed(cloudapiv6.FlagIPVersion) {
		ipVersion := c.Flags().String(cloudapiv6.FlagIPVersion)
		properties.SetIpVersion(ipVersion)

		c.Verbose("Property IP Version set: %v", ipVersion)
	}

	return properties
}

func DeleteAllFirewallRules(c *core.CommandConfig) error {
	datacenterId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	serverId := c.Flags().String(cloudapiv6.ArgServerId)
	nicId := c.Flags().String(cloudapiv6.ArgNicId)

	c.Verbose(constants.DatacenterId, datacenterId)
	c.Verbose("Server ID: %v", serverId)
	c.Verbose("NIC with ID: %v", nicId)
	c.Verbose("Getting Firewall Rules...")
	firewallRules, _, err := c.CloudApiV6Services.FirewallRules().List(datacenterId, serverId, nicId)
	if err != nil {
		return err
	}

	firewallRulesItems, ok := firewallRules.GetItemsOk()
	if !ok || firewallRulesItems == nil {
		return fmt.Errorf("could not get items of Firewall Rules")
	}

	if len(*firewallRulesItems) <= 0 {
		return fmt.Errorf("no Firewall Rule found")
	}

	var multiErr error
	for _, firewall := range *firewallRulesItems {
		id := firewall.GetId()
		name := firewall.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Firewall Rule with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, *id)

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

// checkSourceIPAndTargetIPVersions returns true if the source and destination
// IPs are of a different type (IPv4/IPv6) than the specified IP version.
func checkSourceIPAndTargetIPVersions(c *core.CommandConfig) bool {
	ipVersion := c.Flags().String(cloudapiv6.FlagIPVersion)

	sIp := net.ParseIP(c.Flags().String(cloudapiv6.ArgSourceIp))
	tIp := net.ParseIP(c.Flags().String(cloudapiv6.ArgDestinationIp))

	isIPv4 := func(ip net.IP) bool {
		return ip != nil && ip.To4() != nil
	}

	if (isIPv4(sIp) || isIPv4(tIp)) && ipVersion == "IPv6" {
		return true
	}

	if (!isIPv4(sIp) || !isIPv4(tIp)) && ipVersion == "IPv4" {
		return true
	}

	return false
}
