package flowlog

import (
	"errors"
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

func PreRunApplicationLoadBalancerFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgS3Bucket)
}

func PreRunApplicationLoadBalancerFlowLogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgFlowLogId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func PreRunDcApplicationLoadBalancerFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgFlowLogId)
}

func RunApplicationLoadBalancerFlowLogList(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("Getting FlowLogs")

	applicationloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFlowLogCols).Prefix("items").Print(applicationloadbalancerFlowLogs)
}

func RunApplicationLoadBalancerFlowLogGet(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("Getting FlowLog with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFlowLogCols).Print(ng.FlowLog)
}

func RunApplicationLoadBalancerFlowLogCreate(c *core.CommandConfig) error {
	c.Verbose("Datacenter ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose("ApplicationLoadBalancer ID: %v ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

	proper := helpers.GetFlowLogPropertiesSet(c)
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))

		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if !proper.HasDirection() {
		proper.SetDirection(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))

		c.Verbose("Property Direction set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)))
	}

	if !proper.HasAction() {
		proper.SetAction(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction)))

		c.Verbose("Property Action set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction)))
	}

	c.Verbose("Creating FlowLog")

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.FlowLog{
			FlowLog: ionoscloud.FlowLog{
				Properties: &proper.FlowLogProperties,
			},
		},
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFlowLogCols).Print(ng.FlowLog)
}

func RunApplicationLoadBalancerFlowLogUpdate(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))

	input := helpers.GetFlowLogPropertiesUpdate(c)

	c.Verbose("Updating FlowLog")

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
		&input,
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allFlowLogCols).Print(ng.FlowLog)
}

func RunApplicationLoadBalancerFlowLogDelete(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

		err := DeleteAllApplicationLoadBalancerFlowLog(c)
		if err != nil {
			return err
		}

		return nil
	}

	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("FlowLog ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete application load balancer flowlog with id: %s", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId))), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Application Load Balancers Flowlog successfully deleted")

	return nil
}

func DeleteAllApplicationLoadBalancerFlowLog(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	albId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose(constants.ApplicationLoadBalancerId, albId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.FlowLog]{
		Resource: "Application Load Balancer FlowLog",
		List: func() ([]ionoscloud.FlowLog, error) {
			applicationLoadBalancerFlowlogs, _, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(dcId, albId)
			if err != nil {
				return nil, err
			}
			items, ok := applicationLoadBalancerFlowlogs.GetItemsOk()
			if !ok || items == nil {
				return nil, errors.New("could not get items of Application Load Balancer Flow Logs")
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
			resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(dcId, albId, *fl.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}
