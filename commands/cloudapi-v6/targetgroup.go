package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

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
	globalFlags.StringSliceP(config.ArgCols, "", defaultTargetGroupCols, printer.ColsMessage(allTargetGroupCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(targetGroupCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = targetGroupCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTargetGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
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
	get.AddStringFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Target Group", "Name of the Target Group")
	create.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Protocol of the balancing")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Health Check] It specifies the time (in milliseconds) for a target VM in this pool to answer the check. If a target VM has CheckInterval set and CheckTimeout is set too, then the smaller value of the two is used after the TCP connection is established")
	create.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	create.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 50000, "[Health Check] The maximum inactivity time (in milliseconds) on the target VM side")
	create.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] The number of retries to perform on a target VM after a connection failure. (valid range: [0, 65535])")
	create.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The path for the HTTP health check")
	create.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The method for the HTTP health check")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] Match Type for the HTTP health check")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The response returned by the request")
	create.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Regex for the HTTP health check")
	create.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Negate for the HTTP health check")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Target Group creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Target Group creation [seconds]")

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
	update.AddStringFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Updated Target Group", "Name of the Target Group")
	update.AddStringFlag(cloudapiv6.ArgAlgorithm, "", "ROUND_ROBIN", "Algorithm for the balancing")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAlgorithm, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ROUND_ROBIN", "LEAST_CONNECTION", "RANDOM", "SOURCE_IP"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgProtocol, cloudapiv6.ArgProtocolShort, "HTTP", "Protocol of the balancing")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgProtocol, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HTTP"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddIntFlag(cloudapiv6.ArgCheckTimeout, "", 2000, "[Health Check] It specifies the time (in milliseconds) for a target VM in this pool to answer the check. If a target VM has CheckInterval set and CheckTimeout is set too, then the smaller value of the two is used after the TCP connection is established")
	update.AddIntFlag(cloudapiv6.ArgConnectionTimeout, "", 5000, "[Health Check] It specifies the maximum time (in milliseconds) to wait for a connection attempt to a target VM to succeed")
	update.AddIntFlag(cloudapiv6.ArgTargetTimeout, "", 50000, "[Health Check] The maximum inactivity time (in milliseconds) on the target VM side")
	update.AddIntFlag(cloudapiv6.ArgRetries, "", 3, "[Health Check] The number of retries to perform on a target VM after a connection failure. (valid range: [0, 65535])")
	update.AddStringFlag(cloudapiv6.ArgPath, "", "/.", "[HTTP Health Check] The path for the HTTP health check")
	update.AddStringFlag(cloudapiv6.ArgMethod, "", "GET", "[HTTP Health Check] The method for the HTTP health check")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HEAD", "PUT", "POST", "GET", "TRACE", "PATCH", "OPTIONS"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgMatchType, "", "STATUS_CODE", "[HTTP Health Check] Match Type for the HTTP health check")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"STATUS_CODE", "RESPONSE_BODY"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgResponse, "", "200", "[HTTP Health Check] The response returned by the request")
	update.AddBoolFlag(cloudapiv6.ArgRegex, "", false, "[HTTP Health Check] Regex for the HTTP health check")
	update.AddBoolFlag(cloudapiv6.ArgNegate, "", false, "[HTTP Health Check] Negate for the HTTP health check")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Target Group update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Target Group update [seconds]")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Groups")
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Target Group deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Target Group deletion [seconds]")

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
	c.Printer.Verbose("Getting TargetGroups")
	ss, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(nil, c, getTargetGroups(ss)))
}

func RunTargetGroupGet(c *core.CommandConfig) error {
	c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Printer.Verbose("Getting TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getTargetGroupPrint(nil, c, getTargetGroup(s)))
}

func RunTargetGroupCreate(c *core.CommandConfig) error {
	c.Printer.Verbose("Creating TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Create(getTargetGroupNew(c))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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
	c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	c.Printer.Verbose("Updating TargetGroup")
	s, resp, err := c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), getTargetGroupPropertiesSet(c))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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
	var (
		resp *resources.Response
		err  error
	)
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		resp, err = DeleteAllTargetGroup(c)
		if err != nil {
			return err
		}
	} else {
		c.Printer.Verbose("TargetGroup ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete target group"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting TargetGroup")
		resp, err = c.CloudApiV6Services.TargetGroups().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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

func DeleteAllTargetGroup(c *core.CommandConfig) (*resources.Response, error) {
	_ = c.Printer.Print("Target Groups to be deleted:")
	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if err != nil {
		return nil, err
	}
	if itemsOk, ok := targetGroups.GetItemsOk(); ok && itemsOk != nil {
		for _, alb := range *itemsOk {
			if idOk, ok := alb.GetIdOk(); ok && idOk != nil {
				_ = c.Printer.Print("Target Group Id: " + *idOk)
			}
			if propertiesOk, ok := alb.GetPropertiesOk(); ok && propertiesOk != nil {
				if name, ok := propertiesOk.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print("Target Group Name: " + *name)
				}
			}
		}
		if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Target Groups"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Target Groups...")
		for _, itemOk := range *itemsOk {
			if idOk, ok := itemOk.GetIdOk(); ok && idOk != nil {
				c.Printer.Verbose("Starting deleting Target Group with id: %v...", *idOk)
				resp, err = c.CloudApiV6Services.TargetGroups().Delete(*idOk)
				if resp != nil {
					c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
					c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
				}
				if err != nil {
					return nil, err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, err
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
	//inputHealthCheck.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
	//c.Printer.Verbose("Property ConnectTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
	//inputHealthCheck.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
	//c.Printer.Verbose("Property TargetTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
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

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealthCheck := resources.TargetGroupHealthCheck{}
		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) {
			inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
			c.Printer.Verbose("Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
		}
		//if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)) {
		//	inputHealthCheck.SetConnectTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
		//	c.Printer.Verbose("Property ConnectTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgConnectionTimeout)))
		//}
		//if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)) {
		//	inputHealthCheck.SetTargetTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
		//	c.Printer.Verbose("Property TargetTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgTargetTimeout)))
		//}
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
	defaultTargetGroupCols = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "ConnectTimeout", "TargetTimeout", "State"}
	allTargetGroupCols     = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "ConnectTimeout", "TargetTimeout", "Retries",
		"Path", "Method", "MatchType", "Response", "Regex", "Negate", "State"}
)

type TargetGroupPrint struct {
	TargetGroupId  string `json:"TargetGroupId,omitempty"`
	Name           string `json:"Name,omitempty"`
	Algorithm      string `json:"Algorithm,omitempty"`
	Protocol       string `json:"Protocol,omitempty"`
	CheckTimeout   string `json:"CheckTimeout,omitempty"`
	ConnectTimeout string `json:"ConnectTimeout,omitempty"`
	TargetTimeout  string `json:"TargetTimeout,omitempty"`
	Retries        int32  `json:"Retries,omitempty"`
	Path           string `json:"Path,omitempty"`
	Method         string `json:"Method,omitempty"`
	MatchType      string `json:"MatchType,omitempty"`
	Response       string `json:"Response,omitempty"`
	Regex          bool   `json:"Regex,omitempty"`
	Negate         bool   `json:"Negate,omitempty"`
	State          string `json:"State,omitempty"`
}

func getTargetGroupPrint(resp *resources.Response, c *core.CommandConfig, s []resources.TargetGroup) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getTargetGroupsKVMaps(s)
			r.Columns = getTargetGroupCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getTargetGroupCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultTargetGroupCols
	}
	columnsMap := map[string]string{
		"TargetGroupId":  "TargetGroupId",
		"Name":           "Name",
		"Algorithm":      "Algorithm",
		"Protocol":       "Protocol",
		"CheckTimeout":   "CheckTimeout",
		"ConnectTimeout": "ConnectTimeout",
		"TargetTimeout":  "TargetTimeout",
		"Retries":        "Retries",
		"Path":           "Path",
		"Method":         "Method",
		"MatchType":      "MatchType",
		"Response":       "Response",
		"Regex":          "Regex",
		"Negate":         "Negate",
		"State":          "State",
	}
	var datacenterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			datacenterCols = append(datacenterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
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
			//if connectTimeoutOk, ok := healthCheckOk.GetConnectTimeoutOk(); ok && connectTimeoutOk != nil {
			//	targetGroupPrint.ConnectTimeout = fmt.Sprintf("%dms", *connectTimeoutOk)
			//}
			//if targetTimeoutOk, ok := healthCheckOk.GetTargetTimeoutOk(); ok && targetTimeoutOk != nil {
			//	targetGroupPrint.TargetTimeout = fmt.Sprintf("%dms", *targetTimeoutOk)
			//}
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
