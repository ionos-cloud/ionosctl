package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NlbRuleTargetCmd() *core.Command {
	ctx := context.TODO()
	nlbRuleTargetCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "target",
			Aliases:          []string{"t"},
			Short:            "Network Load Balancer Forwarding Rule Target Operations",
			Long:             "The sub-commands of `ionosctl networkloadbalancer rule target` allow you to add, list, update, remove Network Load Balancer Forwarding Rule Targets.",
			TraverseChildren: true,
		},
	}
	globalFlags := nlbRuleTargetCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultRuleTargetCols, printer.ColsMessage(defaultRuleTargetCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(nlbRuleTargetCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = nlbRuleTargetCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultRuleTargetCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace:  "forwardingrule",
		Resource:   "target",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Network Load Balancer Forwarding Rule Targets",
		LongDesc:   "Use this command to list Targets of a Network Load Balancer Forwarding Rule.\n\nRequired values to run command:\n\n* Data Center Id\n* Network Load Balancer Id\n* Forwarding Rule Id",
		Example:    listNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunDcNetworkLoadBalancerForwardingRuleIds,
		CmdRun:     RunNlbRuleTargetList,
		InitClient: true,
	})
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddUUIDFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	/*
		Add Command
	*/
	add := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "target",
		Verb:      "add",
		Aliases:   []string{"a"},
		ShortDesc: "Add a Network Load Balancer Forwarding Rule Target",
		LongDesc: `Use this command to add a Forwarding Rule Target in a specified Network Load Balancer Forwarding Rule. You can also set Health Check Settings for Forwarding Rule Target. The Check parameter for Health Check Settings specifies whether the target VM's health is checked. If turned off, a target VM is always considered available. If turned on, the target VM is available when accepting periodic TCP connections, to ensure that it is really able to serve requests. The address and port to send the tests to are those of the target VM. The health check only consists of a connection attempt.

Regarding the Weight parameter, this parameter is used to adjust the target VM's weight relative to other target VMs. All target VMs will receive a load proportional to their weight relative to the sum of all weights, so the higher the weight, the higher the load. The default weight is 1, and the maximal value is 256. A value of 0 means the target VM will not participate in load-balancing but will still accept persistent connections. If this parameter is used to distribute the load according to target VM's capacity, it is recommended to start with values which can both grow and shrink, for instance between 10 and 100 to leave enough room above and below for later adjustments.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port`,
		Example:    addNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleTarget,
		CmdRun:     RunNlbRuleTargetAdd,
		InitClient: true,
	})
	add.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddUUIDFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = add.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr,
			viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(add.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	add.AddStringFlag(cloudapiv6.ArgIp, "", "", "IP of a balanced target VM", core.RequiredFlagOption())
	add.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the balanced target service. Range: 1 to 65535", core.RequiredFlagOption())
	add.AddIntFlag(cloudapiv6.ArgWeight, cloudapiv6.ArgWeightShort, 1, "Weight parameter is used to adjust the target VM's weight relative to other target VMs. Maximum: 256")
	add.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] CheckInterval determines the duration (in milliseconds) between consecutive health checks")
	add.AddBoolFlag(cloudapiv6.ArgCheck, "", true, "[Health Check] Check specifies whether the target VM's health is checked")
	add.AddBoolFlag(cloudapiv6.ArgMaintenance, "", false, "[Health Check]  Maintenance specifies if a target VM should be marked as down, even if it is not")
	add.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule Target creation to be executed")
	add.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Target creation [seconds]")
	add.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		Remove Command
	*/
	removeCmd := core.NewCommand(ctx, nlbRuleTargetCmd, core.CommandBuilder{
		Namespace: "forwardingrule",
		Resource:  "target",
		Verb:      "remove",
		Aliases:   []string{"r"},
		ShortDesc: "Remove a Target from a Network Load Balancer Forwarding Rule",
		LongDesc: `Use this command to remove a specified Target from Network Load Balancer Forwarding Rule.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Network Load Balancer Id
* Forwarding Rule Id
* Target Ip
* Target Port`,
		Example:    removeNetworkLoadBalancerRuleTargetExample,
		PreCmdRun:  PreRunNetworkLoadBalancerRuleTargetRemove,
		CmdRun:     RunNlbRuleTargetRemove,
		InitClient: true,
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgNetworkLoadBalancerId, "", "", cloudapiv6.NetworkLoadBalancerId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgNetworkLoadBalancerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NetworkLoadBalancersIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgRuleId, "", "", cloudapiv6.ForwardingRuleId, core.RequiredFlagOption())
	_ = removeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRuleId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ForwardingRulesIds(os.Stderr, viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(removeCmd.NS, cloudapiv6.ArgNetworkLoadBalancerId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeCmd.AddStringFlag(cloudapiv6.ArgIp, "", "", "IP of a balanced target VM", core.RequiredFlagOption())
	removeCmd.AddStringFlag(cloudapiv6.ArgPort, cloudapiv6.ArgPortShort, "", "Port of the balanced target service. Range: 1 to 65535", core.RequiredFlagOption())
	removeCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Forwarding Rule Target deletion to be executed")
	removeCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv6.NlbTimeoutSeconds, "Timeout option for Request for Forwarding Rule Target deletion [seconds]")
	removeCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Forwarding Rule Targets.")
	removeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return nlbRuleTargetCmd
}

func PreRunNetworkLoadBalancerRuleTarget(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgIp, cloudapiv6.ArgPort)
}

func PreRunNetworkLoadBalancerRuleTargetRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgTargetIp, cloudapiv6.ArgTargetPort},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgNetworkLoadBalancerId, cloudapiv6.ArgRuleId, cloudapiv6.ArgAll},
	)
}

func RunNlbRuleTargetList(c *core.CommandConfig) error {
	ng, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		resources.QueryParams{},
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			return c.Printer.Print(getRuleTargetPrint(nil, c, getRuleTargets(targets)))
		} else {
			return errors.New("error getting rule targets")
		}
	} else {
		return errors.New("error getting rule properties")
	}
}

func RunNlbRuleTargetAdd(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	var targetItems []ionoscloud.NetworkLoadBalancerForwardingRuleTarget
	ngOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := ngOld.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			targetItems = *targets
		}
	}
	targetNew := getRuleTargetInfo(c)
	targetItems = append(targetItems, targetNew.NetworkLoadBalancerForwardingRuleTarget)
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))
	nlbForwardingRule := &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}
	c.Printer.Verbose("Adding NlbRuleTarget with id: %v to NetworkLoadBalancer with id: %v", ruleId, nlbId)
	_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbForwardingRule, queryParams)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getRuleTargetPrint(resp, c, []resources.NetworkLoadBalancerForwardingRuleTarget{targetNew}))
}

func RunNlbRuleTargetRemove(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllNlbRuleTarget(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete forwarding rule target"); err != nil {
			return err
		}
		c.Printer.Verbose("NlbRuleTarget with id: %v is removing...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)))
		frOld, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
			queryParams,
		)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		proper, err := getRuleTargetsRemove(c, frOld)
		if err != nil {
			return err
		}
		_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId)),
			proper,
			queryParams,
		)
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getRuleTargetPrint(resp, c, nil))
	}
}

func RemoveAllNlbRuleTarget(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	nlbId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNetworkLoadBalancerId))
	ruleId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRuleId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("NetworkLoadBalancer ID: %v", nlbId)
	c.Printer.Verbose("NetworkLoadBalancerForwardingRule ID: %v", ruleId)
	c.Printer.Verbose("Getting NetworkLoadBalancerForwardingRule...")
	forwardingRule, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().GetForwardingRule(dcId, nlbId, ruleId, cloudapiv6.ParentResourceQueryParams)
	if err != nil {
		return err
	}
	if forwardingRuleProperties, ok := forwardingRule.GetPropertiesOk(); ok && forwardingRuleProperties != nil {
		if targets, ok := forwardingRuleProperties.GetTargetsOk(); ok && targets != nil {
			if len(*targets) > 0 {
				_ = c.Printer.Warn("Forwarding Rule Targets to be removed:")
				for _, target := range *targets {
					toPrint := ""
					if ipOk, ok := target.GetIpOk(); ok && ipOk != nil {
						toPrint += " Forwarding Rule Target IP: " + *ipOk
					}
					if portOk, ok := target.GetPortOk(); ok && portOk != nil {
						toPrint += " Forwarding Rule Target Port: " + strconv.Itoa(int(*portOk))
					}
					_ = c.Printer.Print(toPrint)
				}
			} else {
				return errors.New("no Forwarding Rule Targets found")
			}
		} else {
			return errors.New("could not get items of Forwarding Rule Targets")
		}
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "remove all the Forwarding Rule Targets"); err != nil {
			return err
		}
		c.Printer.Verbose("Removing all the Forwarding Rule Targets...")
		targetItems := make([]ionoscloud.NetworkLoadBalancerForwardingRuleTarget, 0)
		if properties, ok := forwardingRule.GetPropertiesOk(); ok && properties != nil {
			if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
				nlbFwRuleProp := &resources.NetworkLoadBalancerForwardingRuleProperties{
					NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
						Targets: &targetItems,
					},
				}
				_, resp, err = c.CloudApiV6Services.NetworkLoadBalancers().UpdateForwardingRule(dcId, nlbId, ruleId, nlbFwRuleProp, queryParams)
				if resp != nil && printer.GetId(resp) != "" {
					c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
				}
				if err != nil {
					return err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getRuleTargetInfo(c *core.CommandConfig) resources.NetworkLoadBalancerForwardingRuleTarget {
	targetIp := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIp))
	targetPort := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgPort))
	weight := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgWeight))
	target := resources.NetworkLoadBalancerForwardingRuleTarget{}
	target.SetIp(targetIp)
	target.SetPort(targetPort)
	target.SetWeight(weight)
	targetHealth := resources.NetworkLoadBalancerForwardingRuleTargetHealthCheck{}
	maintenance := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgMaintenance))
	check := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCheck))
	checkInterval := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))
	targetHealth.SetMaintenance(maintenance)
	targetHealth.SetCheck(check)
	targetHealth.SetCheckInterval(checkInterval)
	target.SetHealthCheck(targetHealth.NetworkLoadBalancerForwardingRuleTargetHealthCheck)
	c.Printer.Verbose("Properties set for adding the NlbRuleTarget: Ip: %v, Port: %v, Weight: %v, Maintenance: %v, Check: %v, CheckInterval: %v",
		targetIp, targetPort, weight, maintenance, check, checkInterval)
	return target
}

func getRuleTargetsRemove(c *core.CommandConfig, frOld *resources.NetworkLoadBalancerForwardingRule) (*resources.NetworkLoadBalancerForwardingRuleProperties, error) {
	var (
		foundIp   = false
		foundPort = false
	)
	targetItems := make([]ionoscloud.NetworkLoadBalancerForwardingRuleTarget, 0)
	if properties, ok := frOld.GetPropertiesOk(); ok && properties != nil {
		if targets, ok := properties.GetTargetsOk(); ok && targets != nil {
			// Iterate trough all targets
			for _, targetItem := range *targets {
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
	}
	if !foundIp {
		return nil, errors.New("no forwarding rule target with the specified IP found")
	}
	if !foundPort {
		return nil, errors.New("no forwarding rule target with the specified port found")
	}
	return &resources.NetworkLoadBalancerForwardingRuleProperties{
		NetworkLoadBalancerForwardingRuleProperties: ionoscloud.NetworkLoadBalancerForwardingRuleProperties{
			Targets: &targetItems,
		},
	}, nil
}

// Output Printing

var defaultRuleTargetCols = []string{"TargetIp", "TargetPort", "Weight", "Check", "CheckInterval", "Maintenance"}

type RuleTargetPrint struct {
	TargetIp      string `json:"TargetIp,omitempty"`
	TargetPort    int32  `json:"TargetPort,omitempty"`
	Weight        int32  `json:"Weight,omitempty"`
	CheckInterval string `json:"CheckInterval,omitempty"`
	Check         bool   `json:"Check,omitempty"`
	Maintenance   bool   `json:"Maintenance,omitempty"`
}

func getRuleTargetPrint(resp *resources.Response, c *core.CommandConfig, ss []resources.NetworkLoadBalancerForwardingRuleTarget) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if ss != nil {
			r.OutputJSON = ss
			r.KeyValue = getRuleTargetsKVMaps(ss)
			r.Columns = getRuleTargetsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getRuleTargetsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultRuleTargetCols
	}

	columnsMap := map[string]string{
		"TargetIp":      "TargetIp",
		"TargetPort":    "TargetPort",
		"Weight":        "Weight",
		"Check":         "Check",
		"CheckInterval": "CheckInterval",
		"Maintenance":   "Maintenance",
	}
	var ruleTargetCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			ruleTargetCols = append(ruleTargetCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return ruleTargetCols
}

func getRuleTargets(targets *[]ionoscloud.NetworkLoadBalancerForwardingRuleTarget) []resources.NetworkLoadBalancerForwardingRuleTarget {
	ss := make([]resources.NetworkLoadBalancerForwardingRuleTarget, 0)
	if targets != nil {
		for _, s := range *targets {
			ss = append(ss, resources.NetworkLoadBalancerForwardingRuleTarget{
				NetworkLoadBalancerForwardingRuleTarget: s,
			})
		}
	}
	return ss
}

func getRuleTargetsKVMaps(targets []resources.NetworkLoadBalancerForwardingRuleTarget) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(targets))
	for _, target := range targets {
		var targetPrint RuleTargetPrint
		if ip, ok := target.GetIpOk(); ok && ip != nil {
			targetPrint.TargetIp = *ip
		}
		if port, ok := target.GetPortOk(); ok && port != nil {
			targetPrint.TargetPort = *port
		}
		if weight, ok := target.GetWeightOk(); ok && weight != nil {
			targetPrint.Weight = *weight
		}
		if health, ok := target.GetHealthCheckOk(); ok && health != nil {
			if check, ok := health.GetCheckOk(); ok && check != nil {
				targetPrint.Check = *check
			}
			if checkInterval, ok := health.GetCheckIntervalOk(); ok && checkInterval != nil {
				targetPrint.CheckInterval = fmt.Sprintf("%vms", *checkInterval)
			}
			if maintenance, ok := health.GetMaintenanceOk(); ok && maintenance != nil {
				targetPrint.Maintenance = *maintenance
			}
		}
		o := structs.Map(targetPrint)
		out = append(out, o)
	}
	return out
}
