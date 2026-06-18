package flowlog

import (
	"fmt"

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

func PreRunNetworkLoadBalacerFlowLogList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId); err != nil {
		return err
	}
	return nil
}

func PreRunNetworkLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgS3Bucket)
}

func PreRunNetworkLoadBalancerFlowLogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgFlowLogId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func PreRunDcNetworkLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgFlowLogId)
}

func RunNetworkLoadBalancerFlowLogList(c *core.CommandConfig) error {
	networkloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(networkloadbalancerFlowLogs.FlowLogs)
}

func RunNetworkLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	c.Verbose("Network Load Balancer FlowLog with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.FlowLog)
}

func RunNetworkLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	proper := helpers.GetFlowLogPropertiesSet(c)

	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction)))
	}

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.FlowLog)
}

func RunNetworkLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	input := helpers.GetFlowLogPropertiesUpdate(c)

	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		&input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(ng.FlowLog)
}

func RunNetworkLoadBalancerFlowLogDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	networkLoadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNetworkLoadBalancerFlowLogs(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete network load balancer flowlog", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Network Load Balancer FlowLog with id: %v...", flowLogId)

	resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteFlowLog(dcId, networkLoadBalancerId, flowLogId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Network Load Balancer FlowLog successfully deleted")
	return nil
}

func DeleteAllNetworkLoadBalancerFlowLogs(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	networkLoadBalancerId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Network Load Balancer ID: %v", networkLoadBalancerId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.FlowLog]{
		Resource: "Network Load Balancer FlowLog",
		List: func() ([]ionoscloud.FlowLog, error) {
			flowLogs, _, err := c.CloudApiV6Services.NetworkLoadBalancers().ListFlowLogs(dcId, networkLoadBalancerId)
			if err != nil {
				return nil, err
			}
			items, ok := flowLogs.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Network Load Balancer FlowLogs")
			}
			return *items, nil
		},
		Summary: func(fl ionoscloud.FlowLog) string {
			summary := ""
			if props, ok := fl.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil {
					summary += *name
				}
			}
			if id, ok := fl.GetIdOk(); ok && id != nil {
				summary += fmt.Sprintf(" (id: %s)", *id)
			}
			return summary
		},
		ID: func(fl ionoscloud.FlowLog) string {
			if id := fl.GetId(); id != nil {
				return *id
			}
			return ""
		},
		Delete: func(fl ionoscloud.FlowLog) error {
			resp, err := c.CloudApiV6Services.NetworkLoadBalancers().DeleteFlowLog(dcId, networkLoadBalancerId, *fl.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
