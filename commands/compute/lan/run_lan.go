package lan

import (
	"context"
	"errors"
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
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunLanListAll(c)
	}

	lans, resp, err := c.CloudApiV6Services.Lans().List(c.Flags().String(cloudapiv6.ArgDataCenterId))
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
		c.Flags().String(cloudapiv6.ArgLanId), c.Flags().String(cloudapiv6.ArgDataCenterId))

	l, resp, err := c.CloudApiV6Services.Lans().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLanId),
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
	name := c.Flags().String(cloudapiv6.ArgName)
	public := c.Flags().Bool(cloudapiv6.ArgPublic)
	properties := ionoscloud.LanProperties{
		Name:   &name,
		Public: &public,
	}

	c.Verbose("Properties set for creating the Lan: Name: %v, Public: %v", name, public)

	if c.Flags().Changed(cloudapiv6.ArgPccId) {
		pcc := c.Flags().String(cloudapiv6.ArgPccId)
		properties.SetPcc(pcc)

		c.Verbose("Property Pcc set: %v", pcc)
	}

	if c.Flags().Changed(cloudapiv6.FlagIPv6CidrBlock) {
		cidr := strings.ToUpper(c.Flags().String(cloudapiv6.FlagIPv6CidrBlock))

		switch cidr {
		case "DISABLE":
			properties.SetIpv6CidrBlockNil()
		case "AUTO":
			properties.SetIpv6CidrBlock(cidr)
		default:
			cidr = strings.ToLower(cidr)
			dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
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

	c.Verbose("Creating LAN in Datacenter with ID: %v...", c.Flags().String(cloudapiv6.ArgDataCenterId))

	l, resp, err := c.CloudApiV6Services.Lans().Create(c.Flags().String(cloudapiv6.ArgDataCenterId), input)
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

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgPublic) {
		public := c.Flags().Bool(cloudapiv6.ArgPublic)
		input.SetPublic(public)

		c.Verbose("Property Public set: %v", public)
	}

	if c.Flags().Changed(cloudapiv6.ArgPccId) {
		pcc := c.Flags().String(cloudapiv6.ArgPccId)
		input.SetPcc(pcc)

		c.Verbose("Property Pcc set: %v", pcc)
	}

	if c.Flags().Changed(cloudapiv6.FlagIPv6CidrBlock) {
		cidr := strings.ToUpper(c.Flags().String(cloudapiv6.FlagIPv6CidrBlock))

		switch cidr {
		case "DISABLE":
			input.SetIpv6CidrBlockNil()
		case "AUTO":
			input.SetIpv6CidrBlock(cidr)
		default:
			cidr = strings.ToLower(cidr)
			dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
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
		c.Flags().String(cloudapiv6.ArgLanId), c.Flags().String(cloudapiv6.ArgDataCenterId))

	lanUpdated, resp, err := c.CloudApiV6Services.Lans().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLanId),
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
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	lanId := c.Flags().String(cloudapiv6.ArgLanId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
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
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting Lans...")

	lans, resp, err := c.CloudApiV6Services.Lans().List(dcId)
	if err != nil {
		return err
	}

	lansItems, ok := lans.GetItemsOk()
	if !ok || lansItems == nil {
		return fmt.Errorf("could not get items of Lans")
	}

	if len(*lansItems) <= 0 {
		return fmt.Errorf("no Lans found")
	}

	c.Msg("Lans to be deleted:")

	var multiErr error
	for _, lan := range *lansItems {
		id := lan.GetId()
		name := lan.GetProperties().GetName()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Lan with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Lans().Delete(dcId, *id)
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
