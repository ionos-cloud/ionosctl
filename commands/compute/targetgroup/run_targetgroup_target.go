package targetgroup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunTargetGroupIdTargetIpPort(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIp, cloudapiv6.ArgPort)
}

func PreRunTargetGroupTargetRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIp, cloudapiv6.ArgPort},
		[]string{cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgAll},
	)
}

func RunTargetGroupTargetList(c *core.CommandConfig) error {
	c.Verbose(constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Verbose("Getting Targets from TargetGroup")

	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if properties, ok := targetGroups.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			return c.Printer(allTargetGroupTargetCols).Print(*targets)
		} else {
			return errors.New("error getting targets")
		}
	} else {
		return errors.New("error getting properties")
	}
}

func RunTargetGroupTargetAdd(c *core.CommandConfig) error {
	var targetItems []ionoscloud.TargetGroupTarget

	// Get existing Targets from the specified Target Group
	c.Verbose(constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Verbose("Getting TargetGroup")

	targetGroupOld, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if err != nil {
		return err
	}

	if properties, ok := targetGroupOld.GetPropertiesOk(); ok && properties != nil {
		c.Verbose("Getting Targets from TargetGroup")

		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			targetItems = *targets
		}
	}

	targetNew := getTargetGroupTargetInfo(c)

	// Add new Target to the existing Targets in a Target Group
	c.Verbose("Adding new Target to existing Targets")

	targetItems = append(targetItems, targetNew.TargetGroupTarget)

	// Update Target Group with the new Targets
	c.Verbose("Updating TargetGroup with the new Targets")

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{
			TargetGroupProperties: ionoscloud.TargetGroupProperties{
				Targets: &targetItems,
			},
		},
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allTargetGroupTargetCols).Print(targetNew.TargetGroupTarget)
}

func RunTargetGroupTargetRemove(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Verbose(constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))

		_, err := RemoveAllTargetGroupTarget(c)
		if err != nil {
			return err
		}

		return nil
	}

	c.Verbose(constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Verbose("Target IP: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)))
	c.Verbose("Target Port: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort)))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove target from target group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	var propertiesUpdated resources.TargetGroupProperties

	// Get existing Targets from the specified Target Group
	c.Verbose("Getting TargetGroup")

	targetGroupOld, _, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if err != nil {
		return err
	}

	if propertiesOk, ok := targetGroupOld.GetPropertiesOk(); ok && propertiesOk != nil {
		if itemsOk, ok := propertiesOk.GetTargetsOk(); ok && itemsOk != nil {
			// Remove specified Target from Target Group
			c.Verbose("Removing Target from existing Targets")

			newTargets, err := getTargetGroupTargetsRemove(c, itemsOk)
			if err != nil {
				return err
			}

			// Set new Targets for Target Group
			propertiesUpdated.SetTargets(*newTargets)
		}
	}

	// Update Target Group with the new Targets
	c.Verbose("Updating TargetGroup with the new Targets")

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), &propertiesUpdated)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Target Group Target successfully deleted")
	return nil
}

func RemoveAllTargetGroupTarget(c *core.CommandConfig) (*resources.Response, error) {
	c.Msg("Target Group Targets to be deleted:")

	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if err != nil {
		return nil, err
	}

	propertiesOk, ok := applicationLoadBalancerRules.GetPropertiesOk()
	if !ok || propertiesOk == nil {
		return nil, fmt.Errorf("could not retrieve Application Load Balancer properties")
	}

	if httpRulesOk, ok := propertiesOk.GetTargetsOk(); ok && httpRulesOk != nil {
		for _, httpRuleOk := range *httpRulesOk {
			if nameOk, ok := httpRuleOk.GetIpOk(); ok && nameOk != nil {
				c.Msg("Target IP: %v", *nameOk)
			}

			if typeOk, ok := httpRuleOk.GetPortOk(); ok && typeOk != nil {
				c.Msg("Target Port: %v", strconv.Itoa(int(*typeOk)))
			}
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Targets from Target Group", viper.GetBool(constants.ArgForce)) {
		return nil, fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Deleting all the Target Group Targets...")

	propertiesOk.SetTargets([]ionoscloud.TargetGroupTarget{})

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{TargetGroupProperties: *propertiesOk},
	)
	if resp != nil {
		c.Verbose("Request Id: %v", request.GetId(resp))
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return nil, err
	}

	c.Msg("Target Group Targets successfully deleted")
	return resp, err
}

func getTargetGroupTargetInfo(c *core.CommandConfig) resources.TargetGroupTarget {
	target := resources.TargetGroupTarget{}

	target.SetIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)))
	c.Verbose("Property Ip for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)))

	target.SetPort(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort)))
	c.Verbose("Property Port for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort)))

	target.SetWeight(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgWeight)))
	c.Verbose("Property Weight for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgWeight)))

	target.SetMaintenanceEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled)))
	c.Verbose("Property MaintenanceEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled)))

	target.SetHealthCheckEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled)))
	c.Verbose("Property HealthCheckEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled)))

	return target
}

func getTargetGroupTargetsRemove(c *core.CommandConfig, targetsOld *[]ionoscloud.TargetGroupTarget) (*[]ionoscloud.TargetGroupTarget, error) {
	var (
		foundIp   = false
		foundPort = false
	)

	targetItems := make([]ionoscloud.TargetGroupTarget, 0)
	if targetsOld != nil {
		for _, targetItem := range *targetsOld {
			// Iterate trough all targets
			removeIp := false
			removePort := false

			if ip, ok := targetItem.GetIpOk(); ok && ip != nil {
				if *ip == viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)) {
					removeIp = true
					foundIp = true
				}
			}

			if port, ok := targetItem.GetPortOk(); ok && port != nil {
				if *port == viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort)) {
					removePort = true
					foundPort = true
				}
			}

			if removeIp && removePort {
				continue
			} else {
				targetItems = append(targetItems, targetItem)
			}
		}
	}

	if !foundIp {
		return nil, errors.New("no target with the specified IP found")
	}

	if !foundPort {
		return nil, errors.New("no target with the specified port found")
	}
	return &targetItems, nil
}
