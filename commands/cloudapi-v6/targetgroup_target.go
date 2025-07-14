package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
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
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTargetGroupTargetCols = []string{"TargetIp", "TargetPort", "Weight", "HealthCheckEnabled", "MaintenanceEnabled"}
)

func TargetGroupTargetCmd() *core.Command {
	ctx := context.TODO()
	targetGroupTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Target Group Target Operations",
			Long:             "The sub-commands of `ionosctl targetgroup target` allow you to see information, to add, remove Targets from Target Groups.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, targetGroupTargetCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "target",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Target Groups Targets",
		LongDesc:   "Use this command to get a list of Target Groups Targets.",
		Example:    listTargetGroupTargetExample,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupTargetList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = list.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, targetGroupTargetCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Target to a Target Group",
		LongDesc: `Use this command to add a Target to a Target Group. You will need to provide the IP, the port and the weight. Weight parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

Health Check can also be set. The ` + "`" + `--check` + "`" + ` option specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Target Group Id
* Target Ip
* Target Port`,
		Example:    addTargetGroupTargetExample,
		PreCmdRun:  PreRunTargetGroupIdTargetIpPort,
		CmdRun:     RunTargetGroupTargetAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddIpFlag(cloudapiv6.ArgIp, "", nil, "The IP of the balanced target VM.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "The port of the balanced target service; valid range is 1 to 65535.", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "Traffic is distributed in proportion to target weight, relative to the combined weight of all targets. A target with higher weight receives a greater share of traffic. Valid range is 0 to 256 and default is 1; targets with weight of 0 do not participate in load balancing but still accept persistent connections. It is best use values in the middle of the range to leave room for later adjustments.")
	add.AddBoolFlag(cloudapiv6.ArgHealthCheckEnabled, "", true, "Makes the target available only if it accepts periodic health check TCP connection attempts; when turned off, the target is considered always available. The health check only consists of a connection attempt to the address and port of the target. Default is True.")
	add.AddBoolFlag(cloudapiv6.ArgMaintenanceEnabled, cloudapiv6.ArgMaintenanceShort, false, "Maintenance mode prevents the target from receiving balanced traffic.")
	add.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group Target addition to be executed")
	add.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group Target addition [seconds]")
	add.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = add.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Command
	*/
	remove := core.NewCommand(ctx, targetGroupTargetCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "target",
		Verb:       "remove",
		Aliases:    []string{"r"},
		ShortDesc:  "Remove a Target from a Target Group",
		LongDesc:   "Use this command to delete the specified Target from Target Group.\n\nRequired values to run command:\n\n* Target Group Id\n* Target Ip\n* Target Port",
		Example:    removeTargetGroupTargetExample,
		PreCmdRun:  PreRunTargetGroupTargetRemove,
		CmdRun:     RunTargetGroupTargetRemove,
		InitClient: true,
	})
	remove.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = remove.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	remove.AddIpFlag(cloudapiv6.ArgIp, "", nil, "IP of a balanced target VM", core.RequiredFlagOption())
	remove.AddIntFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, 8080, "Port of the balanced target service. (range: 1 to 65535)", core.RequiredFlagOption())
	remove.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Group Targets")
	remove.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group Target deletion to be executed")
	remove.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group Target deletion [seconds]")
	remove.AddStringSliceFlag(constants.ArgCols, "", defaultTargetGroupTargetCols, tabheaders.ColsMessage(defaultTargetGroupTargetCols))
	_ = remove.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultTargetGroupTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(targetGroupTargetCmd, "compute", "")
}

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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Targets from TargetGroup"))

	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), resources.QueryParams{})
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		} else {
			return errors.New("error getting targets")
		}
	} else {
		return errors.New("error getting properties")
	}
}

func RunTargetGroupTargetAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	var targetItems []compute.TargetGroupTarget

	// Get existing Targets from the specified Target Group
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	targetGroupOld, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), queryParams)
	if err != nil {
		return err
	}

	if properties, ok := targetGroupOld.GetPropertiesOk(); ok && properties != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Targets from TargetGroup"))

		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			targetItems = *targets
		}
	}

	targetNew := getTargetGroupTargetInfo(c)

	// Add new Target to the existing Targets in a Target Group
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Adding new Target to existing Targets"))

	targetItems = append(targetItems, targetNew.TargetGroupTarget)

	// Update Target Group with the new Targets
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating TargetGroup with the new Targets"))

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{
			TargetGroupProperties: compute.TargetGroupProperties{
				Targets: &targetItems,
			},
		},
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunTargetGroupTargetRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))

		resp, err = RemoveAllTargetGroupTarget(c)
		if err != nil {
			return err
		}

		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Target IP: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Target Port: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove target from target group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	var propertiesUpdated resources.TargetGroupProperties

	// Get existing Targets from the specified Target Group
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	targetGroupOld, _, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), queryParams)
	if err != nil {
		return err
	}

	if propertiesOk, ok := targetGroupOld.GetPropertiesOk(); ok && propertiesOk != nil {
		if itemsOk, ok := propertiesOk.GetTargetsOk(); ok && itemsOk != nil {
			// Remove specified Target from Target Group
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing Target from existing Targets"))

			newTargets, err := getTargetGroupTargetsRemove(c, itemsOk)
			if err != nil {
				return err
			}

			// Set new Targets for Target Group
			propertiesUpdated.SetTargets(*newTargets)
		}
	}

	// Update Target Group with the new Targets
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating TargetGroup with the new Targets"))

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), &propertiesUpdated, queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Target Group Target successfully deleted"))
	return nil
}

func RemoveAllTargetGroupTarget(c *core.CommandConfig) (*resources.Response, error) {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Target Group Targets to be deleted:"))

	applicationLoadBalancerRules, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), cloudapiv6.ParentResourceQueryParams)
	if err != nil {
		return nil, err
	}

	propertiesOk, ok := applicationLoadBalancerRules.GetPropertiesOk()
	if !ok || propertiesOk == nil {
		return nil, fmt.Errorf("could not retrieve Application Load Balancer properties")
	}

	if httpRulesOk, ok := propertiesOk.GetTargetsOk(); ok && httpRulesOk != nil {
		for _, httpRuleOk := range httpRulesOk {
			if nameOk, ok := httpRuleOk.GetIpOk(); ok && nameOk != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Target IP: %v", *nameOk))
			}

			if typeOk, ok := httpRuleOk.GetPortOk(); ok && typeOk != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Target Port: %v", strconv.Itoa(int(*typeOk))))
			}
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Targets from Target Group", viper.GetBool(constants.ArgForce)) {
		return nil, fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Target Group Targets..."))

	propertiesOk.SetTargets([]compute.TargetGroupTarget{})

	_, resp, err = c.CloudApiV6Services.TargetGroups().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)),
		&resources.TargetGroupProperties{TargetGroupProperties: *propertiesOk},
		resources.QueryParams{},
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Request Id: %v", request.GetId(resp)))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return nil, err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return nil, err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Target Group Targets successfully deleted"))
	return resp, err
}

func getTargetGroupTargetInfo(c *core.CommandConfig) resources.TargetGroupTarget {
	target := resources.TargetGroupTarget{}

	target.SetIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Ip for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))))

	target.SetPort(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Port for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPort))))

	target.SetWeight(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgWeight)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property Weight for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgWeight))))

	target.SetMaintenanceEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property MaintenanceEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenanceEnabled))))

	target.SetHealthCheckEnabled(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Property HealthCheckEnabled for Target set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgHealthCheckEnabled))))

	return target
}

func getTargetGroupTargetsRemove(c *core.CommandConfig, targetsOld *[]compute.TargetGroupTarget) (*[]compute.TargetGroupTarget, error) {
	var (
		foundIp   = false
		foundPort = false
	)

	targetItems := make([]compute.TargetGroupTarget, 0)
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
