package loadbalancer

import (
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

func PreRunLoadBalancerList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunDcLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLoadBalancerId)
}

func PreRunDcLoadBalancerDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgLoadBalancerId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func RunLoadBalancerListAll(c *core.CommandConfig) error {
	datacenters, _, err := c.CloudApiV6Services.DataCenters().List()
	if err != nil {
		return err
	}

	allDcs := helpers.GetDataCenters(datacenters)

	var allLoadbalancers []ionoscloud.Loadbalancers
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		loadBalancers, resp, err := c.CloudApiV6Services.Loadbalancers().List(*id)
		if err != nil {
			return err
		}

		allLoadbalancers = append(allLoadbalancers, loadBalancers.Loadbalancers)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allLoadbalancerCols).Prefix("*.items").Print(allLoadbalancers)
}

func RunLoadBalancerList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunLoadBalancerListAll(c)
	}

	c.Verbose("Getting LoadBalancers from Datacenter with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))

	lbs, resp, err := c.CloudApiV6Services.Loadbalancers().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLoadbalancerCols).Prefix("items").Print(lbs.Loadbalancers)
}

func RunLoadBalancerGet(c *core.CommandConfig) error {
	c.Verbose("Load balancer with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLoadBalancerId)))

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLoadbalancerCols).Print(lb.Loadbalancer)
}

func RunLoadBalancerCreate(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))

	c.Verbose("Properties set for creating the load balancer: Name: %v, Dhcp: %v", name, dhcp)

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Create(dcId, name, dhcp)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLoadbalancerCols).Print(lb.Loadbalancer)
}

func RunLoadBalancerUpdate(c *core.CommandConfig) error {
	input := resources.LoadbalancerProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)

		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgIp)) {
		ip := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
		input.SetIp(ip)

		c.Verbose("Property Ip set: %v", ip)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp)) {
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
		input.SetDhcp(dhcp)

		c.Verbose("Property Dhcp set: %v", dhcp)
	}

	lb, resp, err := c.CloudApiV6Services.Loadbalancers().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLoadBalancerId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allLoadbalancerCols).Print(lb.Loadbalancer)
}

func RunLoadBalancerDelete(c *core.CommandConfig) error {
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	loadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLoadBalancerId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllLoadBalancers(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete loadbalancer", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Load balancer with id: %v is deleting...", loadBalancerId)

	resp, err := c.CloudApiV6Services.Loadbalancers().Delete(dcid, loadBalancerId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Load Balancer successfully deleted")
	return nil
}

func DeleteAllLoadBalancers(c *core.CommandConfig) error {
	dcid := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	c.Verbose(constants.DatacenterId, dcid)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.Loadbalancer]{
		Resource: "Load Balancer",
		List: func() ([]ionoscloud.Loadbalancer, error) {
			loadBalancers, _, err := c.CloudApiV6Services.Loadbalancers().List(dcid)
			if err != nil {
				return nil, err
			}
			items, ok := loadBalancers.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Load Balancers")
			}
			return *items, nil
		},
		Summary: func(lb ionoscloud.Loadbalancer) string {
			summary := ""
			if props, ok := lb.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil {
					summary += *name
				}
				if ip, ok := props.GetIpOk(); ok && ip != nil && *ip != "" {
					summary += fmt.Sprintf(" (public IP: %s)", *ip)
				}
			}
			if id, ok := lb.GetIdOk(); ok && id != nil {
				summary += fmt.Sprintf(" (id: %s)", *id)
			}
			return summary
		},
		ID: func(lb ionoscloud.Loadbalancer) string {
			if id := lb.GetId(); id != nil {
				return *id
			}
			return ""
		},
		Delete: func(lb ionoscloud.Loadbalancer) error {
			resp, err := c.CloudApiV6Services.Loadbalancers().Delete(dcid, *lb.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

func PreRunDataCenterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
}
