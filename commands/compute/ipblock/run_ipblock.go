package ipblock

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

func PreRunIpBlockId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgIpBlockId)
}

func PreRunIpBlockDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgIpBlockId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunIpBlockList(c *core.CommandConfig) error {
	ipblocks, resp, err := c.CloudApiV6Services.IpBlocks().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(ipblocks.IpBlocks)
}

func RunIpBlockGet(c *core.CommandConfig) error {
	c.Verbose("Ip block with id: %v is getting...", c.Flags().String(cloudapiv6.ArgIpBlockId))

	i, resp, err := c.CloudApiV6Services.IpBlocks().Get(c.Flags().String(cloudapiv6.ArgIpBlockId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(i.IpBlock)
}

func RunIpBlockCreate(c *core.CommandConfig) error {
	name := c.Flags().String(cloudapiv6.ArgName)
	loc := c.Flags().String(cloudapiv6.ArgLocation)
	size := c.Flags().Int32(cloudapiv6.ArgSize)

	c.Verbose("Properties set for creating the Ip block: Name: %v, Location: %v, Size: %v", name, loc, size)

	i, resp, err := c.CloudApiV6Services.IpBlocks().Create(name, loc, size)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(i.IpBlock)
}

func RunIpBlockUpdate(c *core.CommandConfig) error {
	input := resources.IpBlockProperties{}
	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	i, resp, err := c.CloudApiV6Services.IpBlocks().Update(c.Flags().String(cloudapiv6.ArgIpBlockId), input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(i.IpBlock)
}

func RunIpBlockDelete(c *core.CommandConfig) error {
	ipBlockId := c.Flags().String(cloudapiv6.ArgIpBlockId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllIpBlocks(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete ipblock", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Ip block with ID: %v...", ipBlockId)

	resp, err := c.CloudApiV6Services.IpBlocks().Delete(ipBlockId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Ip Block successfully deleted")
	return nil
}

func DeleteAllIpBlocks(c *core.CommandConfig) error {
	c.Verbose("Getting all Ip Blocks...")

	ipBlocks, resp, err := c.CloudApiV6Services.IpBlocks().List()
	if err != nil {
		return err
	}

	ipBlocksItems, ok := ipBlocks.GetItemsOk()
	if !ok || ipBlocksItems == nil {
		return fmt.Errorf("could not get items of Ip Blocks")
	}

	if len(*ipBlocksItems) <= 0 {
		return fmt.Errorf("no Ip Blocks found")
	}

	var multiErr error
	for _, dc := range *ipBlocksItems {
		id := dc.GetId()
		name := dc.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the IpBlock with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.IpBlocks().Delete(*id)
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
