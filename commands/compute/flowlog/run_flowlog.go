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

func PreRunFlowLogList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId); err != nil {
		return err
	}
	return nil
}

func PreRunFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgS3Bucket)
}

func PreRunDcServerNicFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgFlowLogId)
}

func PreRunFlowlogDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgFlowLogId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgNicId, cloudapiv6.ArgAll},
	)
}

func RunFlowLogList(c *core.CommandConfig) error {
	flowLogs, resp, err := c.CloudApiV6Services.FlowLogs().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(flowLogs.FlowLogs)
}

func RunFlowLogGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId))

	c.Verbose("FlowLog with id: %v from Nic with id: %v is getting...", flowLogId, nicId)

	flowLog, resp, err := c.CloudApiV6Services.FlowLogs().Get(dcId, serverId, nicId, flowLogId)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(flowLog.FlowLog)
}

func RunFlowLogCreate(c *core.CommandConfig) error {
	properties := helpers.GetFlowLogPropertiesSet(c)
	input := resources.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: &properties.FlowLogProperties,
		},
	}

	flowLog, resp, err := c.CloudApiV6Services.FlowLogs().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(flowLog.FlowLog)
}

func RunFlowLogDelete(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgFlowLogId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllFlowlogs(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete flow log", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting FlowLog with id: %v...", flowLogId)

	resp, err := c.CloudApiV6Services.FlowLogs().Delete(dcId, serverId, nicId, flowLogId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Flowlog successfully deleted")

	return nil

}

func DeleteAllFlowlogs(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNicId))

	c.Verbose(constants.DatacenterId, dcId)
	c.Verbose("Server ID: %v", serverId)
	c.Verbose("NIC ID: %v", nicId)

	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.FlowLog]{
		Resource: "Flowlog",
		List: func() ([]ionoscloud.FlowLog, error) {
			flowlogs, _, err := c.CloudApiV6Services.FlowLogs().List(dcId, serverId, nicId)
			if err != nil {
				return nil, err
			}

			items, ok := flowlogs.GetItemsOk()
			if !ok || items == nil {
				return nil, errors.New("could not get items of Flowlogs")
			}

			return *items, nil
		},
		Summary: func(flowlog ionoscloud.FlowLog) string {
			var id string
			if v, ok := flowlog.GetIdOk(); ok && v != nil {
				id = *v
			}
			summary := fmt.Sprintf("id: %s", id)
			if props, ok := flowlog.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil && *name != "" {
					summary = fmt.Sprintf("%s (name: %s)", summary, *name)
				}
			}
			return summary
		},
		ID: func(flowlog ionoscloud.FlowLog) string {
			if id, ok := flowlog.GetIdOk(); ok && id != nil {
				return *id
			}
			return ""
		},
		Delete: func(flowlog ionoscloud.FlowLog) error {
			resp, err := c.CloudApiV6Services.FlowLogs().Delete(dcId, serverId, nicId, *flowlog.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
