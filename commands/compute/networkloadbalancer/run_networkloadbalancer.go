package networkloadbalancer

import (
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
}

func PreRunNetworkLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunDcNetworkLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId)
}

func PreRunDcNetworkLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunNetworkLoadBalancerListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	var allNetworkLoadBalancers []ionoscloud.NetworkLoadBalancers
	allDcs := helpers.GetDataCenters(datacenters)
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		NetworkLoadBalancers, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().List(*id)
		if err != nil {
			return err
		}

		allNetworkLoadBalancers = append(allNetworkLoadBalancers, NetworkLoadBalancers.NetworkLoadBalancers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allCols).Prefix("*.items").Print(allNetworkLoadBalancers)
}

func RunNetworkLoadBalancerList(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunNetworkLoadBalancerListAll(c)
	}

	networkloadbalancers, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().List(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(networkloadbalancers.NetworkLoadBalancers)
}

func RunNetworkLoadBalancerGet(c *core.CommandConfig) error {
	c.Verbose("Network Load Balancer with id: %v is getting...", c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId))

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Get(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
	)

	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancer)
}

func RunNetworkLoadBalancerCreate(c *core.CommandConfig) error {
	proper := getNewNetworkLoadBalancerInfo(c)

	if !proper.HasName() {
		proper.SetName(c.Flags().String(cloudapiv6.ArgName))
	}

	if !proper.HasTargetLan() {
		proper.SetTargetLan(c.Flags().Int32(cloudapiv6.ArgTargetLan))
	}

	if !proper.HasListenerLan() {
		proper.SetListenerLan(c.Flags().Int32(cloudapiv6.ArgListenerLan))
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Create(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		resources.NetworkLoadBalancer{
			NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
				Properties: &proper.NetworkLoadBalancerProperties,
			},
		},
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancer)
}

func RunNetworkLoadBalancerUpdate(c *core.CommandConfig) error {
	input := getNewNetworkLoadBalancerInfo(c)

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Update(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId),
		*input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.NetworkLoadBalancer)
}

func RunNetworkLoadBalancerDelete(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	nlbId := c.Flags().String(cloudapiv6.ArgNetworkLoadBalancerId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllNetworkLoadBalancers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete network load balancer", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Network Load Balancer with id: %v...", nlbId)

	resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Delete(dcId, nlbId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Network Load Balancer successfully deleted")
	return nil
}

func getNewNetworkLoadBalancerInfo(c *core.CommandConfig) *resources.NetworkLoadBalancerProperties {
	input := ionoscloud.NetworkLoadBalancerProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgIps) {
		ips := c.Flags().StringSlice(cloudapiv6.ArgIps)
		input.SetIps(ips)

		c.Verbose("Property Ips set: %v", ips)
	}

	if c.Flags().Changed(cloudapiv6.ArgListenerLan) {
		listenerLan := c.Flags().Int32(cloudapiv6.ArgListenerLan)
		input.SetListenerLan(listenerLan)

		c.Verbose("Property ListenerLan set: %v", listenerLan)
	}

	if c.Flags().Changed(cloudapiv6.ArgTargetLan) {
		targetLan := c.Flags().Int32(cloudapiv6.ArgTargetLan)
		input.SetTargetLan(targetLan)

		c.Verbose("Property TargetLan set: %v", targetLan)
	}

	if c.Flags().Changed(cloudapiv6.ArgPrivateIps) {
		privateIps := c.Flags().StringSlice(cloudapiv6.ArgPrivateIps)
		input.SetLbPrivateIps(privateIps)

		c.Verbose("Property PrivateIps set: %v", privateIps)
	}

	return &resources.NetworkLoadBalancerProperties{
		NetworkLoadBalancerProperties: input,
	}
}

func DeleteAllNetworkLoadBalancers(c *core.CommandConfig) error {
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Getting Network Load Balancers...")

	networkLoadBalancers, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().List(dcId)
	if err != nil {
		return err
	}

	nlbItems, ok := networkLoadBalancers.GetItemsOk()
	if !ok || nlbItems == nil {
		return fmt.Errorf("could not get items of Network Load Balancers")
	}

	if len(*nlbItems) <= 0 {
		return fmt.Errorf("no Network Load Balancers found")
	}

	var multiErr error
	for _, networkLoadBalancer := range *nlbItems {
		id := networkLoadBalancer.GetId()
		name := networkLoadBalancer.Properties.Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Network Load Balancer with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.NetworkLoadBalancers().Delete(dcId, *id)
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
