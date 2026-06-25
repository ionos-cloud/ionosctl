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
	c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose(constants.ApplicationLoadBalancerId, c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))
	c.Verbose("Getting FlowLogs")

	applicationloadbalancerFlowLogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
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
	c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose(constants.ApplicationLoadBalancerId, c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))
	c.Verbose("Getting FlowLog with ID: %v", c.Flags().String(cloudapiv6.ArgFlowLogId))

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetFlowLog(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgFlowLogId),
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
	c.Verbose("Datacenter ID: %v ", c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose("ApplicationLoadBalancer ID: %v ", c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))

	proper := helpers.GetFlowLogPropertiesSet(c)
	if !proper.HasName() {
		proper.SetName(c.Flags().String(cloudapiv6.ArgName))

		c.Verbose("Property Name set: %v", c.Flags().String(cloudapiv6.ArgName))
	}

	if !proper.HasDirection() {
		proper.SetDirection(c.Flags().String(cloudapiv6.ArgDirection))

		c.Verbose("Property Direction set: %v", c.Flags().String(cloudapiv6.ArgDirection))
	}

	if !proper.HasAction() {
		proper.SetAction(c.Flags().String(cloudapiv6.ArgAction))

		c.Verbose("Property Action set: %v", c.Flags().String(cloudapiv6.ArgAction))
	}

	c.Verbose("Creating FlowLog")

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateFlowLog(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
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
	c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose(constants.ApplicationLoadBalancerId, c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))
	c.Verbose("FlowLog ID: %v", c.Flags().String(cloudapiv6.ArgFlowLogId))

	input := helpers.GetFlowLogPropertiesUpdate(c)

	c.Verbose("Updating FlowLog")

	ng, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateFlowLog(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgFlowLogId),
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

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
		c.Verbose(constants.ApplicationLoadBalancerId, c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))

		err := DeleteAllApplicationLoadBalancerFlowLog(c)
		if err != nil {
			return err
		}

		return nil
	}

	c.Verbose(constants.DatacenterId, c.Flags().String(cloudapiv6.ArgDataCenterId))
	c.Verbose(constants.ApplicationLoadBalancerId, c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId))
	c.Verbose("FlowLog ID: %v", c.Flags().String(cloudapiv6.ArgFlowLogId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete application load balancer flowlog with id: %s", c.Flags().String(cloudapiv6.ArgFlowLogId)), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
		c.Flags().String(cloudapiv6.ArgFlowLogId),
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
	c.Msg("Getting Application Load Balancer FlowLogs...")

	applicationLoadBalancerFlowlogs, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListFlowLogs(
		c.Flags().String(cloudapiv6.ArgDataCenterId),
		c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId),
	)
	if err != nil {
		return err
	}

	albFlowLogItems, ok := applicationLoadBalancerFlowlogs.GetItemsOk()
	if !ok || albFlowLogItems == nil {
		return errors.New("could not get items of Application Load Balancer Flow Logs")
	}

	if len(*albFlowLogItems) <= 0 {
		return errors.New("no Application Load Balancer Flow Logs found")
	}

	var multiErr error
	for _, fl := range *albFlowLogItems {
		id := fl.GetId()
		name := fl.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Application Load Balancer FlowLog Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteFlowLog(
			c.Flags().String(cloudapiv6.ArgDataCenterId),
			c.Flags().String(cloudapiv6.ArgApplicationLoadBalancerId), *id,
		)
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

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}
