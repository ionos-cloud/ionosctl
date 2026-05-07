package rule

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

func PreRunDcApplicationLoadBalancerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId)
}

func PreRunApplicationLoadBalancerForwardingRuleDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgAll},
	)
}

func PreRunApplicationLoadBalancerForwardingRuleCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgListenerIp, cloudapiv6.ArgListenerPort)
}

func PreRunDcApplicationLoadBalancerForwardingRuleIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgApplicationLoadBalancerId, cloudapiv6.ArgRuleId)
}

func RunApplicationLoadBalancerForwardingRuleList(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("Getting ForwardingRules")

	albForwardingRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allAlbForwardingRuleCols).Prefix("items").Print(albForwardingRules.ApplicationLoadBalancerForwardingRules)
}

func RunApplicationLoadBalancerForwardingRuleGet(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose("Getting ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allAlbForwardingRuleCols).Print(applicationLoadBalancerForwardingRule.ApplicationLoadBalancerForwardingRule)
}

func RunApplicationLoadBalancerForwardingRuleCreate(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

	proper := getAlbForwardingRulePropertiesSet(c)
	if !proper.HasProtocol() {
		proper.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}
	if !proper.HasName() {
		proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	c.Verbose("Creating ForwardingRule")

	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().CreateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		resources.ApplicationLoadBalancerForwardingRule{
			ApplicationLoadBalancerForwardingRule: ionoscloud.ApplicationLoadBalancerForwardingRule{
				Properties: &proper.ApplicationLoadBalancerForwardingRuleProperties,
			},
		},
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allAlbForwardingRuleCols).Print(applicationLoadBalancerForwardingRule.ApplicationLoadBalancerForwardingRule)
}

func RunApplicationLoadBalancerForwardingRuleUpdate(c *core.CommandConfig) error {
	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	input := getAlbForwardingRulePropertiesSet(c)

	c.Verbose("Updating ForwardingRule")

	applicationLoadBalancerForwardingRule, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().UpdateForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		input,
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allAlbForwardingRuleCols).Print(applicationLoadBalancerForwardingRule.ApplicationLoadBalancerForwardingRule)
}

func RunApplicationLoadBalancerForwardingRuleDelete(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
		c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))

		err := DeleteAllApplicationLoadBalancerForwardingRule(c)
		if err != nil {
			return err
		}

		return nil
	}

	c.Verbose(constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Verbose(constants.ApplicationLoadBalancerId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)))
	c.Verbose(constants.ForwardingRuleId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete application load balancer forwarding rule", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Deleting ForwardingRule with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))

	resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().DeleteForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Application Load Balancer Forwarding Rule successfully deleted")

	return nil
}

func DeleteAllApplicationLoadBalancerForwardingRule(c *core.CommandConfig) error {
	c.Msg("Getting Application Load Balancer Forwarding Rules...")

	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().ListForwardingRules(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)),
	)
	if err != nil {
		return err
	}

	albRuleItems, ok := applicationLoadBalancerRules.GetItemsOk()
	if !ok || albRuleItems == nil {
		return errors.New("could not get items of Target Groups")
	}

	if len(*albRuleItems) <= 0 {
		return errors.New("no Target Groups found")
	}

	var multiErr error
	for _, fr := range *albRuleItems {
		id := fr.GetId()
		name := fr.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete Forwarding Rule Id: %s , Name: %s ", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.ApplicationLoadBalancers().DeleteForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgApplicationLoadBalancerId)), *id,
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

func getAlbForwardingRulePropertiesSet(c *core.CommandConfig) *resources.ApplicationLoadBalancerForwardingRuleProperties {
	input := ionoscloud.ApplicationLoadBalancerForwardingRuleProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)) {
		input.SetListenerIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)))
		c.Verbose("Property ListenerIp set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgListenerIp)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)) {
		input.SetListenerPort(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)))
		c.Verbose("Property ListenerPort set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgListenerPort)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)) {
		input.SetServerCertificates(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
		c.Verbose("Property ServerCertificates set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgServerCertificates)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)) {
		input.SetClientTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
		c.Verbose("Property Client Timeout set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgClientTimeout)))
	}

	return &resources.ApplicationLoadBalancerForwardingRuleProperties{
		ApplicationLoadBalancerForwardingRuleProperties: input,
	}
}
