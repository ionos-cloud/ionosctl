package nic

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunNicList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId); err != nil {
		return err
	}
	return nil
}

func PreRunNicDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgAll},
	)
}

func PreRunNicCreate(c *core.PreCommandConfig) error {
	if err := PreRunDcServerIds(c); err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6IPs)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) {
		return fmt.Errorf("IPv6 IPs cannot be explicitly set unless a Cidr Block is also specified")
	}

	return nil
}

func RunNicList(c *core.CommandConfig) error {

	nics, resp, err := c.CloudApiV6Services.Nics().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Prefix("items").Print(nics.Nics)
}

func RunNicGet(c *core.CommandConfig) error {
	c.Verbose("Nic with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)))

	n, resp, err := c.CloudApiV6Services.Nics().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Print(n.Nic)
}

func RunNicCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	ips := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
	lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
	firewallActive := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallActive))
	firewallType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallType))

	c.Verbose("Creating Nic in DataCenterId: %v with ServerId: %v...", dcId, serverId)
	c.Verbose("Properties set for creating the Nic: Name: %v, Ips: %v, Dhcp: %v, Lan: %v FirewallActive: %v, FirewallType: %v",
		name, ips, dhcp, lanId, firewallActive, firewallType)

	inputProper := resources.NicProperties{}
	inputProper.SetName(name)
	inputProper.SetIps(ips)
	inputProper.SetDhcp(dhcp)
	inputProper.SetLan(lanId)
	inputProper.SetFirewallActive(firewallActive)
	inputProper.SetFirewallType(firewallType)

	lan, resp, err := c.CloudApiV6Services.Lans().Get(dcId, fmt.Sprintf("%d", lanId))
	if err != nil && resp != nil && resp.StatusCode != 404 {
		// Only non-404 errors are returned.
		// If the LAN does not exist, it will be created when the NIC is created.
		return err
	}
	// If LAN exists, check if IPv6 is enabled
	if err == nil {
		isIPv6, err := checkIPv6EnableForLAN(lan.Lan)
		if err != nil {
			return err
		}

		if isIPv6 {
			if err = setIPv6Properties(c, &inputProper.NicProperties, lan.Lan); err != nil {
				return err
			}
		} else if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDHCPv6)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6IPs)) {
			return fmt.Errorf("IPv6 is not enabled on the LAN that the NIC is on")
		}
	}

	input := resources.Nic{
		Nic: ionoscloud.Nic{
			Properties: &inputProper.NicProperties,
		},
	}

	n, resp, err := c.CloudApiV6Services.Nics().Create(dcId, serverId, input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Print(n.Nic)
}

func RunNicUpdate(c *core.CommandConfig) error {
	input := getNicProperties(c)
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	svId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	oldNIc, _, err := c.CloudApiV6Services.Nics().Get(dcId, svId, nicId)
	if err != nil {
		return err
	}

	lan, _, err := c.CloudApiV6Services.Lans().Get(dcId, fmt.Sprintf("%d", *oldNIc.Properties.Lan))
	if err != nil {
		return err
	}

	isIPv6, err := checkIPv6EnableForLAN(lan.Lan)
	if err != nil {
		return err
	}

	if isIPv6 {
		input.NicProperties.SetIpv6CidrBlock(*oldNIc.Properties.Ipv6CidrBlock)
		if err = setIPv6Properties(c, &input.NicProperties, lan.Lan); err != nil {
			return err
		}
	} else if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagDHCPv6)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6IPs)) {
		return fmt.Errorf("IPv6 is not enabled on the LAN that the NIC is on")
	}

	nicUpd, resp, err := c.CloudApiV6Services.Nics().Update(dcId, svId, nicId, input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Print(nicUpd.Nic)
}

func RunNicDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNics(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nic", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Nic with id: %v...", nicId)

	resp, err := c.CloudApiV6Services.Nics().Delete(dcId, serverId, nicId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Nic successfully deleted")
	return nil
}

func DeleteAllNics(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Server ID: %v", serverId)
	c.Verbose("Getting NICs...")

	nics, resp, err := c.CloudApiV6Services.Nics().List(dcId, serverId)
	if err != nil {
		return err
	}

	nicsItems, ok := nics.GetItemsOk()
	if !ok || nicsItems == nil {
		return fmt.Errorf("could not get items of NICs")
	}

	if len(*nicsItems) <= 0 {
		return fmt.Errorf("no NICs found")
	}

	var multiErr error
	for _, nic := range *nicsItems {
		id := nic.GetId()
		name := nic.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Nic with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Nics().Delete(dcId, serverId, *id)
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

func validateIPv6IPs(cidr string, ips ...string) error {
	_, ipNet, _ := net.ParseCIDR(cidr)

	for _, ipString := range ips {
		ip := net.ParseIP(ipString)
		if ip == nil {
			return fmt.Errorf("failed parsing \"%s\" as an IP", ipString)
		}

		if ip.To4() != nil {
			return fmt.Errorf("\"%s\" is not an IPv6 IP", ipString)
		}

		if !ipNet.Contains(ip) {
			return fmt.Errorf("the provided IPv6 IP (%s) is not within the NIC IPv6 Cidr Block", ip)
		}
	}
	return nil
}

func checkIPv6EnableForLAN(lan ionoscloud.Lan) (bool, error) {
	cidr, err := helpers.GetIPv6CidrBlockFromLAN(lan)
	if err != nil {
		return false, err
	}
	if cidr == "" {
		return false, nil
	}

	return true, nil
}

func setIPv6Properties(c *core.CommandConfig, inputProper *ionoscloud.NicProperties, lan ionoscloud.Lan) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) {
		cidr := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)))
		lanIPv6CidrBlock, err := helpers.GetIPv6CidrBlockFromLAN(lan)
		if err != nil {
			return err
		}

		if err := utils2.ValidateIPv6CidrBlockAgainstParentCidrBlock(cidr, 80, lanIPv6CidrBlock); err != nil {
			return err
		}

		inputProper.SetIpv6CidrBlock(cidr)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6IPs)) && inputProper.Ipv6CidrBlock != nil {
		ipv6Ips, _ := c.Command.Command.Flags().GetStringSlice(cloudapiv6.FlagIPv6IPs)
		cidr := *inputProper.Ipv6CidrBlock
		if err := validateIPv6IPs(cidr, ipv6Ips...); err != nil {
			return err
		}

		inputProper.SetIpv6Ips(ipv6Ips)
	}

	dhcpv6 := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.FlagDHCPv6))
	inputProper.SetDhcpv6(dhcpv6)

	return nil
}

func getNicProperties(c *core.CommandConfig) resources.NicProperties {
	input := resources.NicProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.NicProperties.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
		input.NicProperties.SetDhcp(dhcp)

		c.Verbose("Property Dhcp set: %v", dhcp)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
		lan := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
		input.NicProperties.SetLan(lan)

		c.Verbose("Property Lan set: %v", lan)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		ips := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps))
		input.NicProperties.SetIps(ips)

		c.Verbose("Property Ips set: %v", ips)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallActive)) {
		firewallActive := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallActive))
		input.NicProperties.SetFirewallActive(firewallActive)

		c.Verbose("Property FirewallActive set: %v", firewallActive)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallType)) {
		firewallType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFirewallType))
		input.NicProperties.SetFirewallType(firewallType)

		c.Verbose("Property FirewallType set: %v", firewallType)
	}

	return input
}

func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)
}

func PreRunDcServerNicIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId)
}
