package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TargetGroupCmd() *core.Command {
	ctx := context.TODO()
	targetGroupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "targetgroup",
			Aliases:          []string{"tg"},
			Short:            "Target Group Operations",
			Long:             "The sub-commands of `ionosctl targetgroup` allow you to see information, to create, update, delete Target Groups.",
			TraverseChildren: true,
		},
	}
	globalFlags := targetGroupCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTargetGroupCols, printer.ColsMessage(allTargetGroupCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(targetGroupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = targetGroupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTargetGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "targetgroup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Target Groups",
		LongDesc:   "Use this command to get a list of Target Groups.",
		Example:    listTargetGroupExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTargetGroupList,
		InitClient: true,
	})
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "targetgroup",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Target Group",
		LongDesc:   "Use this command to get information about a specified Target Group.\n\nRequired values to run command:\n\n* Target Group Id",
		Example:    getTargetGroupExample,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "targetgroup",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Target Group",
		LongDesc: `Use this command to create a Target Group.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.`,
		Example:    createTargetGroupExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTargetGroupCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Target Group", "The name of the target group.")
	create.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Balancing algorithm.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Balancing protocol")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Health Check] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.")
	create.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] The interval in milliseconds between consecutive health checks; default is 2000.")
	create.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection attempts.")
	create.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The path (destination URL) for the HTTP health check request; the default is /.")
	create.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The method for the HTTP health check.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] Match Type for the HTTP health check.")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMatchType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The response returned by the request, depending on the match type.")
	create.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Regex for the HTTP health check.")
	create.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Negate for the HTTP health check.")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group creation to be executed.")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group creation [seconds].")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "targetgroup",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Target Group",
		LongDesc: `Use this command to update a specified Target Group.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` or ` + "`" + `-w` + "`" + ` option.

Required values to run command:

* Target Group Id`,
		Example:    updateTargetGroupExample,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Updated Target Group", "The name of the target group.")
	update.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Balancing algorithm.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Balancing protocol")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Health Check] The maximum time in milliseconds to wait for a target to respond to a check. For target VMs with 'Check Interval' set, the lesser of the two  values is used once the TCP connection is established.")
	update.AddIntFlag(cloudapiv6.ArgCheckInterval, "", 2000, "[Health Check] The interval in milliseconds between consecutive health checks; default is 2000.")
	update.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] The maximum number of attempts to reconnect to a target after a connection failure. Valid range is 0 to 65535, and default is three reconnection attempts.")
	update.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The path (destination URL) for the HTTP health check request; the default is /.")
	update.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The method for the HTTP health check.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] Match Type for the HTTP health check.")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The response returned by the request, depending on the match type.")
	update.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Regex for the HTTP health check.")
	update.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Negate for the HTTP health check.")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group update to be executed.")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group update [seconds].")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "targetgroup",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete a Target Group",
		LongDesc:   "Use this command to delete the specified Target Group.\n\nRequired values to run command:\n\n* Target Group Id",
		Example:    deleteTargetGroupExample,
		PreCmdRun:  PreRunTargetGroupDelete,
		CmdRun:     RunTargetGroupDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Groups")
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group deletion to be executed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	targetGroupCmd.AddCommand(TargetGroupTargetCmd())

	return targetGroupCmd
}

func PreRunTargetGroupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgTargetGroupId)
}

func PreRunTargetGroupDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgTargetGroupId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunTargetGroupList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	c.Printer.Verbose("Getting TargetGroups")
	ss, resp, err := c.CloudApiV6Services.TargetGroups().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(nil, c, getTargetGroups(ss)))
}

func RunTargetGroupGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Printer.Verbose("Getting TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(nil, c, getTargetGroup(s)))
}

func RunTargetGroupCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Creating TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Create(getTargetGroupNew(c), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(resp, c, getTargetGroup(s)))
}

func RunTargetGroupUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Printer.Verbose("Updating TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), getTargetGroupPropertiesSet(c), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(resp, c, getTargetGroup(s)))
}

func RunTargetGroupDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	var (
		resp *resources.Response
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		err = DeleteAllTargetGroup(c)
		if err != nil {
			return err
		}
	} else {
		c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete target group"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting TargetGroup")
		resp, err = c.CloudApiV6Services.TargetGroups().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), queryParams)
		if resp != nil {
			c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getTargetGroupPrint(resp, c, nil))
}

func DeleteAllTargetGroup(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	_ = c.Printer.Warn("Getting Target Groups...")
	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if targetGroupItems, ok := targetGroups.GetItemsOk(); ok && targetGroupItems != nil {
		if len(*targetGroupItems) > 0 {
			for _, tg := range *targetGroupItems {
				delIdAndName := ""
				if id, ok := tg.GetIdOk(); ok && id != nil {
					delIdAndName += "Target Group Id: " + *id
				}
				if properties, ok := tg.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += "Target Group Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Target Groups"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Target Groups...")
			var multiErr error
			for _, tg := range *targetGroupItems {
				if id, ok := tg.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Target Group with id: %v...", *id)
					resp, err = c.CloudApiV6Services.TargetGroups().Delete(*id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Target Groups found")
		}
	} else {
		return errors.New("could not get items of Target Groups")
	}
}

func getTargetGroupNew(c *core.CommandConfig) resources.TargetGroup {
	input := resources.TargetGroupProperties{}
	// Set Required Properties
	input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	c.Printer.Verbose("Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))

	inputHealthCheck := resources.TargetGroupHealthCheck{}
	// Set Properties for Health Check for Target Group
	inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
	c.Printer.Verbose("Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
	inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
	c.Printer.Verbose("Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
	inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	c.Printer.Verbose("Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	// Set Health Check for Target Group
	input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
	c.Printer.Verbose("Setting HealthCheck")

	inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}
	// Set Properties for Http Health Check for Target Group
	inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
	c.Printer.Verbose("Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
	inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
	c.Printer.Verbose("Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
	inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
	c.Printer.Verbose("Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
	inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	c.Printer.Verbose("Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
	c.Printer.Verbose("Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
	inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	c.Printer.Verbose("Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	// Set Http Health Check for Target Group
	input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
	c.Printer.Verbose("Setting HttpHealthCheck")

	return resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Properties: &input.TargetGroupProperties,
		},
	}
}

func getTargetGroupPropertiesSet(c *core.CommandConfig) *resources.TargetGroupProperties {
	input := resources.TargetGroupProperties{}
	// Set new values for Required Properties
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
		c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)) {
		input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
		c.Printer.Verbose("Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
		c.Printer.Verbose("Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealthCheck := resources.TargetGroupHealthCheck{}
		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) {
			inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
			c.Printer.Verbose("Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) {
			inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
			c.Printer.Verbose("Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
			inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
			c.Printer.Verbose("Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
		}
		input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
		c.Printer.Verbose("Updating HealthCheck")
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
		inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}
		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) {
			inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
			c.Printer.Verbose("Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) {
			inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
			c.Printer.Verbose("Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) {
			inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
			c.Printer.Verbose("Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
			inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
			c.Printer.Verbose("Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) {
			inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
			c.Printer.Verbose("Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) {
			inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
			c.Printer.Verbose("Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
		}
		input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
		c.Printer.Verbose("Updating HttpHealthCheck")
	}

	return &input
}

// Output Printing

var (
	defaultTargetGroupCols = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "State"}
	allTargetGroupCols     = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "Retries",
		"Path", "Method", "MatchType", "Response", "Regex", "Negate", "State"}
)

type TargetGroupPrint struct {
	TargetGroupId string `json:"TargetGroupId,omitempty"`
	Name          string `json:"Name,omitempty"`
	Algorithm     string `json:"Algorithm,omitempty"`
	Protocol      string `json:"Protocol,omitempty"`
	CheckTimeout  string `json:"CheckTimeout,omitempty"`
	CheckInterval string `json:"CheckInterval,omitempty"`
	Retries       int32  `json:"Retries,omitempty"`
	Path          string `json:"Path,omitempty"`
	Method        string `json:"Method,omitempty"`
	MatchType     string `json:"MatchType,omitempty"`
	Response      string `json:"Response,omitempty"`
	Regex         bool   `json:"Regex,omitempty"`
	Negate        bool   `json:"Negate,omitempty"`
	State         string `json:"State,omitempty"`
}

func getTargetGroupPrint(resp *resources.Response, c *core.CommandConfig, s []resources.TargetGroup) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getTargetGroupsKVMaps(s)
			r.Columns = printer.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols)))
		}
	}
	return r
}

func getTargetGroups(targetGroups resources.TargetGroups) []resources.TargetGroup {
	ss := make([]resources.TargetGroup, 0)
	if items, ok := targetGroups.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.TargetGroup{TargetGroup: s})
		}
	}
	return ss
}

func getTargetGroup(s *resources.TargetGroup) []resources.TargetGroup {
	ss := make([]resources.TargetGroup, 0)
	if s != nil {
		ss = append(ss, resources.TargetGroup{TargetGroup: s.TargetGroup})
	}
	return ss
}

func getTargetGroupsKVMaps(ss []resources.TargetGroup) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getTargetGroupKVMap(s)
		out = append(out, o)
	}
	return out
}

func getTargetGroupKVMap(s resources.TargetGroup) map[string]interface{} {
	var targetGroupPrint TargetGroupPrint
	if ssId, ok := s.GetIdOk(); ok && ssId != nil {
		targetGroupPrint.TargetGroupId = *ssId
	}
	if propertiesOk, ok := s.GetPropertiesOk(); ok && propertiesOk != nil {
		if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
			targetGroupPrint.Name = *nameOk
		}
		if algorithmOk, ok := propertiesOk.GetAlgorithmOk(); ok && algorithmOk != nil {
			targetGroupPrint.Algorithm = *algorithmOk
		}
		if protocolOk, ok := propertiesOk.GetProtocolOk(); ok && protocolOk != nil {
			targetGroupPrint.Protocol = *protocolOk
		}
		if healthCheckOk, ok := propertiesOk.GetHealthCheckOk(); ok && healthCheckOk != nil {
			if checkTimeoutOk, ok := healthCheckOk.GetCheckTimeoutOk(); ok && checkTimeoutOk != nil {
				targetGroupPrint.CheckTimeout = fmt.Sprintf("%dms", *checkTimeoutOk)
			}
			if checkIntervalOk, ok := healthCheckOk.GetCheckIntervalOk(); ok && checkIntervalOk != nil {
				targetGroupPrint.CheckInterval = fmt.Sprintf("%dms", *checkIntervalOk)
			}
			if retriesOk, ok := healthCheckOk.GetRetriesOk(); ok && retriesOk != nil {
				targetGroupPrint.Retries = *retriesOk
			}
		}
		if httpHealthCheckOk, ok := propertiesOk.GetHttpHealthCheckOk(); ok && httpHealthCheckOk != nil {
			if pathOk, ok := httpHealthCheckOk.GetPathOk(); ok && pathOk != nil {
				targetGroupPrint.Path = *pathOk
			}
			if methodOk, ok := httpHealthCheckOk.GetMethodOk(); ok && methodOk != nil {
				targetGroupPrint.Method = *methodOk
			}
			if matchTypeOk, ok := httpHealthCheckOk.GetMatchTypeOk(); ok && matchTypeOk != nil {
				targetGroupPrint.MatchType = *matchTypeOk
			}
			if responseOk, ok := httpHealthCheckOk.GetResponseOk(); ok && responseOk != nil {
				targetGroupPrint.Response = *responseOk
			}
			if regexOk, ok := httpHealthCheckOk.GetRegexOk(); ok && regexOk != nil {
				targetGroupPrint.Regex = *regexOk
			}
			if negateOk, ok := httpHealthCheckOk.GetNegateOk(); ok && negateOk != nil {
				targetGroupPrint.Negate = *negateOk
			}
		}
	}
	if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			targetGroupPrint.State = *state
		}
	}
	return structs.Map(targetGroupPrint)
}
