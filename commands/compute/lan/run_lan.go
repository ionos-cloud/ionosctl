package lan

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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

func PreRunLansList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunLanDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunLanListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)

	var allLans []ionoscloud.Lans
	totalTime := time.Duration(0)
	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("failed to retrieve Datacenter ID")
		}

		lans, resp, err := c.CloudApiV6Services.Lans().List(*dc.GetId())
		if err != nil {
			return err
		}

		allLans = append(allLans, lans.Lans)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allLanCols).Prefix("*.items").Print(allLans)
}

func RunLanList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunLanListAll(c)
	}

	lans, resp, err := c.CloudApiV6Services.Lans().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLanCols).Prefix("items").Print(lans.Lans)
}

func RunLanGet(c *core.CommandConfig) error {
	c.Verbose("Lan with id: %v from Datacenter with id: %v is getting...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	l, resp, err := c.CloudApiV6Services.Lans().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLanCols).Print(l.Lan)
}

func RunLanCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
	properties := ionoscloud.LanProperties{
		Name:   &name,
		Public: &public,
	}

	c.Verbose("Properties set for creating the Lan: Name: %v, Public: %v", name, public)

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		properties.SetPcc(pcc)

		c.Verbose("Property Pcc set: %v", pcc)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) {
		cidr := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)))

		switch cidr {
		case "DISABLE":
			properties.SetIpv6CidrBlockNil()
		case "AUTO":
			properties.SetIpv6CidrBlock(cidr)
		default:
			cidr = strings.ToLower(cidr)
			dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
			dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
			if err != nil {
				return err
			}

			dcIPv6CidrBlock, err := helpers.GetIPv6CidrBlockFromDatacenter(dc)
			if err != nil {
				return err
			}

			if err = utils2.ValidateIPv6CidrBlockAgainstParentCidrBlock(cidr, 64, dcIPv6CidrBlock); err != nil {
				return err
			}

			properties.SetIpv6CidrBlock(cidr)
		}

		c.Verbose("Property IPv6 Cidr Block set: %v", cidr)
	}

	input := resources.LanPost{
		Lan: ionoscloud.Lan{
			Properties: &properties,
		},
	}

	c.Verbose("Creating LAN in Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	l, resp, err := c.CloudApiV6Services.Lans().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLanCols).Print(l.Lan)
}

func RunLanUpdate(c *core.CommandConfig) error {
	input := resources.LanProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublic)) {
		public := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic))
		input.SetPublic(public)

		c.Verbose("Property Public set: %v", public)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)) {
		pcc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))
		input.SetPcc(pcc)

		c.Verbose("Property Pcc set: %v", pcc)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)) {
		cidr := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.FlagIPv6CidrBlock)))

		switch cidr {
		case "DISABLE":
			input.SetIpv6CidrBlockNil()
		case "AUTO":
			input.SetIpv6CidrBlock(cidr)
		default:
			cidr = strings.ToLower(cidr)
			dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
			dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
			if err != nil {
				return err
			}

			dcIPv6CidrBlock, err := helpers.GetIPv6CidrBlockFromDatacenter(dc)
			if err != nil {
				return err
			}

			if err = utils2.ValidateIPv6CidrBlockAgainstParentCidrBlock(cidr, 64, dcIPv6CidrBlock); err != nil {
				return err
			}

			input.SetIpv6CidrBlock(cidr)
		}
	}

	c.Verbose("Updating LAN with ID: %v from Datacenter with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)), viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLanCols).Print(lanUpdated.Lan)
}

func RunLanDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	lanId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllLans(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete lan", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting LAN with ID: %v from Datacenter with ID: %v...", lanId, dcId)

	resp, err := c.CloudApiV6Services.Lans().Delete(dcId, lanId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Lan successfully deleted")
	return nil
}

func DeleteAllLans(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	c.Verbose(constants.DatacenterId, dcId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.Lan]{
		Resource: "lan",
		List: func() ([]ionoscloud.Lan, error) {
			lans, _, err := c.CloudApiV6Services.Lans().List(dcId)
			if err != nil {
				return nil, err
			}

			items, ok := lans.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Lans")
			}

			return *items, nil
		},
		Summary: func(lan ionoscloud.Lan) string {
			var id, name string
			if lan.Id != nil {
				id = *lan.Id
			}
			if p := lan.Properties; p != nil && p.Name != nil {
				name = *p.Name
			}
			return fmt.Sprintf("%s (id: %s)", name, id)
		},
		ID: func(lan ionoscloud.Lan) string {
			if lan.Id != nil {
				return *lan.Id
			}
			return ""
		},
		Delete: func(lan ionoscloud.Lan) error {
			resp, err := c.CloudApiV6Services.Lans().Delete(dcId, *lan.Id)
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

func GetIPv6CidrBlockFromLAN(lan ionoscloud.Lan) (string, error) {
	if properties, ok := lan.GetPropertiesOk(); ok && properties != nil {
		if ipv6CidrBlock, ok := properties.GetIpv6CidrBlockOk(); ok && ipv6CidrBlock != nil {
			return *ipv6CidrBlock, nil
		} else if ok && ipv6CidrBlock == nil {
			return "", nil
		}
	}

	return "", fmt.Errorf("could not retrieve IPv6 Cidr Block from LAN")
}

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
}

func PreRunDcLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLanId)
}
