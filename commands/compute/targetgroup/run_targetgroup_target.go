package targetgroup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Targets from TargetGroup"))

	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if properties, ok := targetGroups.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroupTarget, *targets,
				tabheaders.GetHeadersAllDefault(defaultTargetGroupTargetCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	targetGroupOld, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if err != nil {
		return err
	}

	if properties, ok := targetGroupOld.GetPropertiesOk(); ok && properties != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Targets from TargetGroup"))

		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			targetItems = *targets
		}
	}

	targetNew := getTargetGroupTargetInfo(c)

	// Add new Target to the existing Targets in a Target Group
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Adding new Target to existing Targets"))

	targetItems = append(targetItems, targetNew.TargetGroupTarget)

	// Update Target Group with the new Targets
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating TargetGroup with the new Targets"))

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{
			TargetGroupProperties: ionoscloud.TargetGroupProperties{
				Targets: &targetItems,
			},
		},
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroupTarget, targetNew.TargetGroupTarget,
		tabheaders.GetHeadersAllDefault(defaultTargetGroupTargetCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupTargetRemove(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))

		_, err := RemoveAllTargetGroupTarget(c)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Target IP: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Target Port: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove target from target group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	var propertiesUpdated resources.TargetGroupProperties

	// Get existing Targets from the specified Target Group
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	targetGroupOld, _, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if err != nil {
		return err
	}

	if propertiesOk, ok := targetGroupOld.GetPropertiesOk(); ok && propertiesOk != nil {
		if itemsOk, ok := propertiesOk.GetTargetsOk(); ok && itemsOk != nil {
			// Remove specified Target from Target Group
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing Target from existing Targets"))

			newTargets, err := getTargetGroupTargetsRemove(c, itemsOk)
			if err != nil {
				return err
			}

			// Set new Targets for Target Group
			propertiesUpdated.SetTargets(*newTargets)
		}
	}

	// Update Target Group with the new Targets
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating TargetGroup with the new Targets"))

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), &propertiesUpdated)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Target Group Target successfully deleted"))
	return nil
}

func RemoveAllTargetGroupTarget(c *core.CommandConfig) (*resources.Response, error) {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Target Group Targets to be deleted:"))

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
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Target IP: %v", *nameOk))
			}

			if typeOk, ok := httpRuleOk.GetPortOk(); ok && typeOk != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Target Port: %v", strconv.Itoa(int(*typeOk))))
			}
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Targets from Target Group", viper.GetBool(constants.ArgForce)) {
		return nil, fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all the Target Group Targets..."))

	propertiesOk.SetTargets([]ionoscloud.TargetGroupTarget{})

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{TargetGroupProperties: *propertiesOk},
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Request Id: %v", request.GetId(resp)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return nil, err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return nil, err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Target Group Targets successfully deleted"))
	return resp, err
}

func getTargetGroupTargetInfo(c *core.CommandConfig) resources.TargetGroupTarget {
	target := resources.TargetGroupTarget{}

	target.SetIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Ip for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))))

	target.SetPort(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Port for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort))))

	target.SetWeight(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgWeight)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Weight for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgWeight))))

	target.SetMaintenanceEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property MaintenanceEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled))))

	target.SetHealthCheckEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property HealthCheckEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled))))

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
