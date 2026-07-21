package pcc

import (
	"errors"
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
	c.Verbose("Cross Connect with id: %v is getting...", c.Flags().String(cloudapiv6.ArgPccId))

	u, resp, err := c.CloudApiV6Services.Pccs().Get(c.Flags().String(cloudapiv6.ArgPccId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Print(u.PrivateCrossConnect)
}

func RunPccCreate(c *core.CommandConfig) error {
	name := c.Flags().String(cloudapiv6.ArgName)
	description := c.Flags().String(cloudapiv6.ArgDescription)

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
	oldPcc, resp, err := c.CloudApiV6Services.Pccs().Get(c.Flags().String(cloudapiv6.ArgPccId))
	if err != nil {
		return err
	}

	newProperties := getPccInfo(oldPcc, c)
	pccUpd, resp, err := c.CloudApiV6Services.Pccs().Update(c.Flags().String(cloudapiv6.ArgPccId), *newProperties)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allPccCols).Print(pccUpd.PrivateCrossConnect)
}

func RunPccDelete(c *core.CommandConfig) error {
	pccId := c.Flags().String(cloudapiv6.ArgPccId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
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
		if c.Flags().Changed(cloudapiv6.ArgName) {
			namePcc = c.Flags().String(cloudapiv6.ArgName)

			c.Verbose("Property Name set: %v", namePcc)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				namePcc = *name
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgDescription) {
			description = c.Flags().String(cloudapiv6.ArgDescription)

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
	c.Verbose("Getting PrivateCrossConnects...")

	pccs, resp, err := c.CloudApiV6Services.Pccs().List()
	if err != nil {
		return err
	}

	pccsItems, ok := pccs.GetItemsOk()
	if !ok || pccsItems == nil {
		return fmt.Errorf("could not get items of PrivateCrossConnects")
	}

	if len(*pccsItems) <= 0 {
		return fmt.Errorf("no PrivateCrossConnects found")
	}

	var multiErr error
	for _, pcc := range *pccsItems {
		id := pcc.GetId()
		name := pcc.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the PrivateCrossConnect with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Pccs().Delete(*id)
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

func RunPccPeersList(c *core.CommandConfig) error {
	c.Verbose("Getting Peers from Cross-Connect with ID: %v...", c.Flags().String(cloudapiv6.ArgPccId))

	u, resp, err := c.CloudApiV6Services.Pccs().GetPeers(c.Flags().String(cloudapiv6.ArgPccId))
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
