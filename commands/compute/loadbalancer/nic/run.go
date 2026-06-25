package nic

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunLoadBalancerNicList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLoadBalancerId); err != nil {
		return err
	}
	return nil
}

func PreRunDcNicLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNicId, cloudapiv6.ArgLoadBalancerId)
}

func PreRunNicDetach(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNicId, cloudapiv6.ArgLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNicId, cloudapiv6.ArgAll},
	)
}

func RunLoadBalancerNicAttach(c *core.CommandConfig) error {
	attachedNic, _, err := c.CloudApiV6Services.Loadbalancers().AttachNic(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgNicId),
	)
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Print(attachedNic.Nic)
}

func RunLoadBalancerNicList(c *core.CommandConfig) error {

	attachedNics, _, err := c.CloudApiV6Services.Loadbalancers().ListNics(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLoadBalancerId),
	)
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Prefix("items").Print(attachedNics.BalancedNics)
}

func RunLoadBalancerNicGet(c *core.CommandConfig) error {
	n, _, err := c.CloudApiV6Services.Loadbalancers().GetNic(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgNicId),
	)
	if err != nil {
		return err
	}

	return c.Printer(allNicCols).Print(n.Nic)
}

func RunLoadBalancerNicDetach(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DetachAllNics(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "detach nic from loadbalancer", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := c.CloudApiV6Services.Loadbalancers().DetachNic(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgNicId),
	)
	if err != nil {
		return err
	}

	c.Msg("Nic successfully detached from Load Balancer")
	return nil
}

func DetachAllNics(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	lbId := c.Flags().String(cloudapiv6.ArgLoadBalancerId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("LoadBalancer ID: %v", lbId)
	c.Verbose("Getting NICs...")

	nics, resp, err := c.CloudApiV6Services.Loadbalancers().ListNics(dcId, lbId)
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

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Detach the Nic with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Loadbalancers().DetachNic(dcId, lbId, *id)
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
