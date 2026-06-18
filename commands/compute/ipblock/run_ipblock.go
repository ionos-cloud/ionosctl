package ipblock

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
	c.Verbose("Ip block with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))

	i, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(i.IpBlock)
}

func RunIpBlockCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	loc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	size := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgSize))

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
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	i, resp, err := c.CloudApiV6Services.IpBlocks().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)), input)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(i.IpBlock)
}

func RunIpBlockDelete(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
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
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.IpBlock]{
		Resource: "ipblock",
		List: func() ([]ionoscloud.IpBlock, error) {
			ipBlocks, _, err := c.CloudApiV6Services.IpBlocks().List()
			if err != nil {
				return nil, err
			}

			items, ok := ipBlocks.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Ip Blocks")
			}

			return *items, nil
		},
		Summary: func(ipBlock ionoscloud.IpBlock) string {
			var id, name, location string
			var ips []string
			if ipBlock.Id != nil {
				id = *ipBlock.Id
			}
			if p := ipBlock.Properties; p != nil {
				if p.Name != nil {
					name = *p.Name
				}
				if p.Location != nil {
					location = *p.Location
				}
				if p.Ips != nil {
					ips = *p.Ips
				}
			}
			if len(ips) > 0 {
				return fmt.Sprintf("%s (id: %s, ips: %s, location: %s)", name, id, strings.Join(ips, ", "), location)
			}
			return fmt.Sprintf("%s (id: %s, location: %s)", name, id, location)
		},
		ID: func(ipBlock ionoscloud.IpBlock) string {
			if ipBlock.Id != nil {
				return *ipBlock.Id
			}
			return ""
		},
		Delete: func(ipBlock ionoscloud.IpBlock) error {
			resp, err := c.CloudApiV6Services.IpBlocks().Delete(*ipBlock.Id)
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
