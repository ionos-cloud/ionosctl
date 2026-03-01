package firewallrule

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
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
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.FirewallRule, firewallRules.FirewallRules,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunFirewallRuleGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Getting Firewall Rule with id: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId))))

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunFirewallRuleCreate(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	properties := getFirewallRulePropertiesSet(c)

	if !properties.HasName() {
		properties.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if !properties.HasType() {
		properties.SetType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))
	}

	input := resources.FirewallRule{
		FirewallRule: ionoscloud.FirewallRule{
			Properties: &properties.FirewallruleProperties,
		},
	}
	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunFirewallRuleUpdate(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		if checkSourceIPAndTargetIPVersions(c) {
			return fmt.Errorf("if source IP and destination IP are set, they must be the same version as IP version")
		}
	}

	firewallRule, resp, err := c.CloudApiV6Services.FirewallRules().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId)),
		getFirewallRulePropertiesSet(c),
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.FirewallRule, firewallRule.FirewallRule,
		tabheaders.GetHeaders(allFirewallRuleCols, defaultFirewallRuleCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunFirewallRuleDelete(c *core.CommandConfig) error {
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))
	fruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallRuleId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Firewall Rule successfully deleted"))

	return nil

}

// Get Firewall Rule Properties set used for create and update commands
func getFirewallRulePropertiesSet(c *core.CommandConfig) resources.FirewallRuleProperties {
	properties := resources.FirewallRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		properties.SetName(name)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		protocol := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))
		properties.SetProtocol(protocol)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Protocol set: %v", protocol))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceIp)) {
		sourceIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceIp))
		properties.SetSourceIp(sourceIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property SourceIp set: %v", sourceIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSourceMac)) {
		sourceMac := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceMac))
		properties.SetSourceMac(sourceMac)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property SourceMac set: %v", sourceMac))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDestinationIp)) {
		targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDestinationIp))
		properties.SetTargetIp(targetIp)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property TargetIp/DestinationIp set: %v", targetIp))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpCode)) {
		icmpCode := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpCode))
		properties.SetIcmpCode(icmpCode)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property IcmpCode set: %v", icmpCode))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpType)) {
		icmpType := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgIcmpType))
		properties.SetIcmpType(icmpType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property IcmpType set: %v", icmpType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart)) {
		portRangeStart := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeStart))
		properties.SetPortRangeStart(portRangeStart)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property PortRangeStart set: %v", portRangeStart))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd)) {
		portRangeEnd := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPortRangeEnd))
		properties.SetPortRangeEnd(portRangeEnd)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property PortRangeEnd set: %v", portRangeEnd))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)) {
		firewallruleType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection))
		properties.SetType(strings.ToUpper(firewallruleType))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Type/Direction set: %v", firewallruleType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion)) {
		ipVersion := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion))
		properties.SetIpVersion(ipVersion)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property IP Version set: %v", ipVersion))
	}

	return properties
}

func DeleteAllFirewallRules(c *core.CommandConfig) error {
	datacenterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, datacenterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Server ID: %v", serverId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("NIC with ID: %v", nicId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Firewall Rules..."))
	firewallRules, resp, err := c.CloudApiV6Services.FirewallRules().List(datacenterId, serverId, nicId)
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

		resp, err = c.CloudApiV6Services.FirewallRules().Delete(datacenterId, serverId, nicId, *id)

		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
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
	ipVersion := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPVersion))

	sIp := net.ParseIP(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSourceIp)))
	tIp := net.ParseIP(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDestinationIp)))

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
