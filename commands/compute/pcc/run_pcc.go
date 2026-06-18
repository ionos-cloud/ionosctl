package pcc

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

func PreRunPccId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgPccId)
}

func PreRunPccDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgPccId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunPccList(c *core.CommandConfig) error {

	pccs, resp, err := c.CloudApiV6Services.Pccs().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Prefix("items").Print(pccs.PrivateCrossConnects)
}

func RunPccGet(c *core.CommandConfig) error {
	c.Verbose("Cross Connect with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))

	u, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Print(u.PrivateCrossConnect)
}

func RunPccCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	description := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))

	newUser := resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Properties: &ionoscloud.PrivateCrossConnectProperties{
				Name:        &name,
				Description: &description,
			},
		},
	}

	c.Verbose("Properties set for creating the Cross Connect: Name: %v, Description: %v", name, description)

	u, resp, err := c.CloudApiV6Services.Pccs().Create(newUser)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Print(u.PrivateCrossConnect)
}

func RunPccUpdate(c *core.CommandConfig) error {
	oldPcc, resp, err := c.CloudApiV6Services.Pccs().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if err != nil {
		return err
	}

	newProperties := getPccInfo(oldPcc, c)
	pccUpd, resp, err := c.CloudApiV6Services.Pccs().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)), *newProperties)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Print(pccUpd.PrivateCrossConnect)
}

func RunPccDelete(c *core.CommandConfig) error {
	pccId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllPccs(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete Cross-Connect", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Cross Connect with id: %v...", pccId)

	resp, err := c.CloudApiV6Services.Pccs().Delete(pccId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Cross Connect successfully deleted")
	return nil
}

func getPccInfo(oldUser *resources.PrivateCrossConnect, c *core.CommandConfig) *resources.PrivateCrossConnectProperties {
	var namePcc, description string

	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
			namePcc = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))

			c.Verbose("Property Name set: %v", namePcc)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				namePcc = *name
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDescription)) {
			description = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDescription))

			c.Verbose("Property Description set: %v", description)
		} else {
			if desc, ok := properties.GetDescriptionOk(); ok && desc != nil {
				description = *desc
			}
		}
	}

	return &resources.PrivateCrossConnectProperties{
		PrivateCrossConnectProperties: ionoscloud.PrivateCrossConnectProperties{
			Name:        &namePcc,
			Description: &description,
		},
	}
}

func DeleteAllPccs(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.PrivateCrossConnect]{
		Resource: "PrivateCrossConnect",
		List: func() ([]ionoscloud.PrivateCrossConnect, error) {
			pccs, _, err := c.CloudApiV6Services.Pccs().List()
			if err != nil {
				return nil, err
			}

			items, ok := pccs.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of PrivateCrossConnects")
			}

			return *items, nil
		},
		Summary: func(pcc ionoscloud.PrivateCrossConnect) string {
			summary := fmt.Sprintf("id: %s", *pcc.GetId())
			if props, ok := pcc.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil && *name != "" {
					summary = fmt.Sprintf("%s (name: %s)", summary, *name)
				}
				if desc, ok := props.GetDescriptionOk(); ok && desc != nil && *desc != "" {
					summary = fmt.Sprintf("%s (desc: %s)", summary, *desc)
				}
			}
			return summary
		},
		ID: func(pcc ionoscloud.PrivateCrossConnect) string {
			return *pcc.GetId()
		},
		Delete: func(pcc ionoscloud.PrivateCrossConnect) error {
			resp, err := c.CloudApiV6Services.Pccs().Delete(*pcc.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

func RunPccPeersList(c *core.CommandConfig) error {
	c.Verbose("Getting Peers from Cross-Connect with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))

	u, resp, err := c.CloudApiV6Services.Pccs().GetPeers(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPccId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	peers := make([]ionoscloud.Peer, 0)

	if u != nil {
		for _, p := range *u {
			peers = append(peers, p.Peer)
		}
	}

	return c.Printer(allPccPeerCols).Print(peers)
}
