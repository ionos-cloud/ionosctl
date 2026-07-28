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

func PreRunNATGatewayFlowLogList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId); err != nil {
		return err
	}
	return nil
}

func PreRunNatGatewayFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgS3Bucket)
}

func PreRunNatGatewayFlowlogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgFlowLogId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgAll},
	)
}

func PreRunDcNatGatewayFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNatGatewayId, cloudapiv6.ArgFlowLogId)
}

func RunNatGatewayFlowLogList(c *core.CommandConfig) error {
	natgatewayFlowLogs, resp, err := c.CloudApiV6Services.NatGateways().ListFlowLogs(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(natgatewayFlowLogs.FlowLogs)
}

func RunNatGatewayFlowLogGet(c *core.CommandConfig) error {
	c.Verbose("NatGatewayFlowLogGet with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId)))

	ng, resp, err := c.CloudApiV6Services.NatGateways().GetFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
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

func RunNatGatewayFlowLogCreate(c *core.CommandConfig) error {
	proper := helpers.GetFlowLogPropertiesSet(c)

	ng, resp, err := c.CloudApiV6Services.NatGateways().CreateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
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

func RunNatGatewayFlowLogUpdate(c *core.CommandConfig) error {
	input := helpers.GetFlowLogPropertiesUpdate(c)

	ng, resp, err := c.CloudApiV6Services.NatGateways().UpdateFlowLog(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId)),
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

func RunNatGatewayFlowLogDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natgatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))
	flowlogId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNatGatewayFlowLogs(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete nat gateway flowlog", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting NatGatewayFlowLog with id: %v...", flowlogId)

	resp, err := c.CloudApiV6Services.NatGateways().DeleteFlowLog(dcId, natgatewayId, flowlogId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("NAT Gateway Flowlog successfully deleted")
	return nil
}

func DeleteAllNatGatewayFlowLogs(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	natgatewayId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("NatGateway ID: %v", natgatewayId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.FlowLog]{
		Resource: "NAT Gateway FlowLog",
		List: func() ([]ionoscloud.FlowLog, error) {
			flowlogs, _, err := c.CloudApiV6Services.NatGateways().ListFlowLogs(dcId, natgatewayId)
			if err != nil {
				return nil, err
			}
			items, ok := flowlogs.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of NAT Gateway FlowLogs")
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
			resp, err := c.CloudApiV6Services.NatGateways().DeleteFlowLog(dcId, natgatewayId, *fl.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
