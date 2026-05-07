package applicationloadbalancer

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

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}

func PreRunApplicationLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func PreRunApplicationLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func RunApplicationLoadBalancerListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)
	var allApplicationLoadBalancers = make([]ionoscloud.ApplicationLoadBalancers, 0)
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		dcId, ok := dc.GetIdOk()
		if !ok || dcId == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		ApplicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(*dcId)
		if err != nil {
			return err
		}

		allApplicationLoadBalancers = append(allApplicationLoadBalancers, ApplicationLoadBalancers.ApplicationLoadBalancers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allApplicationLoadBalancerCols).Prefix("*.items").Print(allApplicationLoadBalancers)
}

func RunApplicationLoadBalancerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunApplicationLoadBalancerListAll(c)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	c.Verbose("Getting ApplicationLoadBalancers from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	applicationloadbalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(dcId)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	if err != nil {
		return err
	}

	return c.Printer(allApplicationLoadBalancerCols).Prefix("items").Print(applicationloadbalancers)
}

func RunApplicationLoadBalancerGet(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose("Getting ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	if err != nil {
		return err
	}

	return c.Printer(allApplicationLoadBalancerCols).Print(ng.ApplicationLoadBalancer)
}

func RunApplicationLoadBalancerCreate(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	proper := getNewApplicationLoadBalancerInfo(c)

	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if !proper.HasTargetLan() {
		proper.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Verbose("Property TargetLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}

	if !proper.HasListenerLan() {
		proper.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Verbose("Property ListenerLan set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}

	c.Verbose("Creating ApplicationLoadBalancer")
	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		resources.ApplicationLoadBalancer{
			ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
				Properties: &proper.ApplicationLoadBalancerProperties,
			},
		},
	)
	if resp != nil {
		c.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allApplicationLoadBalancerCols).Print(ng.ApplicationLoadBalancer)
}

func RunApplicationLoadBalancerUpdate(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	input := getNewApplicationLoadBalancerInfo(c)
	c.Verbose("Updating ApplicationLoadBalancer with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		*input,
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allApplicationLoadBalancerCols).Print(ng.ApplicationLoadBalancer)
}

func RunApplicationLoadBalancerDelete(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

		err := DeleteAllApplicationLoadBalancer(c)
		if err != nil {
			return err
		}

		return nil
	}

	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete application load balancer", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting ApplicationLoadBalancer")

	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Application Load Balancer successfully deleted")

	return nil
}

func DeleteAllApplicationLoadBalancer(c *core.CommandConfig) error {
	c.Msg("Getting Application Load Balancers...")

	applicationLoadBalancers, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if err != nil {
		return err
	}

	albItems, ok := applicationLoadBalancers.GetItemsOk()
	if !ok || albItems == nil {
		return errors.New("could not get items of Application Load Balancers")
	}

	if len(*albItems) <= 0 {
		return errors.New("no Application Load Balancers found")
	}

	var multiErr error

	for _, alb := range *albItems {
		id := alb.GetId()
		name := alb.Properties.Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Application Load Balancer Id: %s , Name: %s ", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *id)
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

func getNewApplicationLoadBalancerInfo(c *core.CommandConfig) *resources.ApplicationLoadBalancerProperties {
	input := ionoscloud.ApplicationLoadBalancerProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIps)) {
		input.SetIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
		c.Verbose("Property IPs set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgIps)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)) {
		input.SetListenerLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
		c.Verbose("Property ListenerLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)) {
		input.SetTargetLan(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
		c.Verbose("Property TargetLan set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetLan)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)) {
		input.SetLbPrivateIps(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
		c.Verbose("Property LbPrivateIps set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPrivateIps)))
	}
	return &resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: input,
	}
}
