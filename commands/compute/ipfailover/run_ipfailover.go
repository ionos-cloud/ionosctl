package ipfailover

import (
	"fmt"

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

func PreRunDcLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId)
}

func PreRunDcLanServerNicIpRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgIp},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgAll},
	)
}

func PreRunDcLanServerNicIdsIp(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgIp)
}

func RunIpFailoverList(c *core.CommandConfig) error {
	ipsFailovers := make([]ionoscloud.IPFailover, 0)
	obj, resp, err := c.CloudApiV6Services.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := obj.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting lan properties")
	}

	ipFailovers, ok := properties.GetIpFailoverOk()
	if !ok || ipFailovers == nil {
		return fmt.Errorf("error getting ip failovers")
	}

	for _, ip := range *ipFailovers {
		ipsFailovers = append(ipsFailovers, ip)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.IpFailover, ipsFailovers,
		tabheaders.GetHeadersAllDefault(defaultIpFailoverCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunIpFailoverAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Adding an IP Failover group to LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId))

	ipsFailovers := make([]ionoscloud.IPFailover, 0)
	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(dcId, lanId, getIpFailoverInfo(c))
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	properties, ok := lanUpdated.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting updated lan properties")
	}

	ipFailovers, ok := properties.GetIpFailoverOk()
	if !ok || ipFailovers == nil {
		return fmt.Errorf("error getting updated ipfailovers")
	}

	for _, ip := range *ipFailovers {
		ipsFailovers = append(ipsFailovers, ip)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.IpFailover, ipsFailovers,
		tabheaders.GetHeadersAllDefault(defaultIpFailoverCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunIpFailoverRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllIpFailovers(c); err != nil {
			return err
		}

		return nil
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Removing IP Failover group from LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove ip failover group from lan", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	oldLan, _, err := c.CloudApiV6Services.Lans().Get(dcId, lanId)
	if err != nil {
		return err
	}

	properties, ok := oldLan.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting lan properties to update")
	}

	ipfailovers, ok := properties.GetIpFailoverOk()
	if !ok || ipfailovers == nil {
		return fmt.Errorf("error getting ipfailovers to update")
	}

	_, resp, err := c.CloudApiV6Services.Lans().Update(dcId, lanId, removeIpFailoverInfo(c, ipfailovers))
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Ip Failover successfully deleted"))
	return nil
}

func RemoveAllIpFailovers(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))

	newIpFailover := make([]ionoscloud.IPFailover, 0)
	lanProperties := resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &newIpFailover,
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Lan ID: %v", lanId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing IP Failovers..."))

	ipFailovers, resp, err := c.CloudApiV6Services.Lans().List(dcId)
	if err != nil {
		return err
	}

	ipFailoversItems, ok := ipFailovers.GetItemsOk()
	if !ok || ipFailoversItems == nil {
		return fmt.Errorf("could not get items of IP Failovers")
	}

	if len(*ipFailoversItems) <= 0 {
		return fmt.Errorf("no IP Failovers found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("IP Failovers to be removed:"))

	for _, ipFailover := range *ipFailoversItems {
		delIdAndName := ""
		if id, ok := ipFailover.GetIdOk(); ok && id != nil {
			delIdAndName += "IP Failover Id: " + *id
		}

		if properties, ok := ipFailover.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " IP Failover Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("%s", delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the IP Failovers", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	oldLan, _, err := c.CloudApiV6Services.Lans().Get(dcId, lanId)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing all the IP Failovers..."))

	if properties, ok := oldLan.GetPropertiesOk(); ok && properties != nil {
		if ipfailovers, ok := properties.GetIpFailoverOk(); ok && ipfailovers != nil {
			_, resp, err = c.CloudApiV6Services.Lans().Update(dcId, lanId, lanProperties)
			if resp != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Request Id: %v", request.GetId(resp)))
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
			}
			if err != nil {
				return err
			}

			if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
				return err
			}
		}
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Ip Failovers successfully deleted"))
	return nil
}

func getIpFailoverInfo(c *core.CommandConfig) resources.LanProperties {
	ip := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Adding IpFailover with Ip: %v and NicUuid: %v", ip, nicId))

	return resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &[]ionoscloud.IPFailover{
				{
					Ip:      &ip,
					NicUuid: &nicId,
				},
			},
		},
	}
}

func removeIpFailoverInfo(c *core.CommandConfig, failovers *[]ionoscloud.IPFailover) resources.LanProperties {
	removeIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	removeNicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing IpFailover with Ip: %v and NicUuid: %v", removeIp, removeNicId))

	newIpFailover := make([]ionoscloud.IPFailover, 0)
	for _, failover := range *failovers {
		if ip, ok := failover.GetIpOk(); ok && ip != nil && *ip != removeIp {
			if nicId, ok := failover.GetNicUuidOk(); ok && nicId != nil && *nicId != removeNicId {
				newIpFailover = append(newIpFailover, failover)
			}
		}
	}

	return resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &newIpFailover,
		},
	}
}
