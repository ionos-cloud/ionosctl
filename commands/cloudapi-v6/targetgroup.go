package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
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
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultTargetGroupCols = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "State"}
	allTargetGroupCols     = []string{"TargetGroupId", "Name", "Algorithm", "Protocol", "CheckTimeout", "CheckInterval", "Retries",
		"Path", "Method", "MatchType", "Response", "Regex", "Negate", "State"}
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultTargetGroupCols, tabheaders.ColsMessage(allTargetGroupCols))
	_ = viper.BindPFlag(core.GetFlagName(targetGroupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = targetGroupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allTargetGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	_ = core.NewCommand(ctx, targetGroupCmd, core.CommandBuilder{
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
	get.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
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
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
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
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Groups")
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Target Group deletion to be executed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Target Group deletion [seconds]")

	targetGroupCmd.AddCommand(TargetGroupTargetCmd())

	return core.WithConfigOverride(targetGroupCmd, []string{fileconfiguration.Cloud, "compute"}, "")
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroups"))
	ss, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.TargetGroup, ss.TargetGroups,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupCreate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Create(getTargetGroupNew(c))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupUpdate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating TargetGroup"))

	s, resp, err := c.CloudApiV6Services.TargetGroups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)), getTargetGroupPropertiesSet(c))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.TargetGroup, s.TargetGroup,
		tabheaders.GetHeaders(allTargetGroupCols, defaultTargetGroupCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunTargetGroupDelete(c *core.CommandConfig) error {
	var resp *resources.Response

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))
		err := DeleteAllTargetGroup(c)
		if err != nil {
			return err
		}

		return nil
	}
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		constants.TargetGroupId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId))))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete target group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting TargetGroup"))

	resp, err := c.CloudApiV6Services.TargetGroups().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgTargetGroupId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Target Group successfully deleted"))
	return nil
}

func DeleteAllTargetGroup(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Getting Target Groups..."))

	targetGroups, resp, err := c.CloudApiV6Services.TargetGroups().List()
	if err != nil {
		return err
	}

	targetGroupItems, ok := targetGroups.GetItemsOk()
	if !ok || targetGroupItems == nil {
		return fmt.Errorf("could not get items of Target Groups")
	}

	if len(*targetGroupItems) <= 0 {
		return fmt.Errorf("no Target Groups found")
	}

	var multiErr error
	for _, tg := range *targetGroupItems {
		id := tg.GetId()
		name := tg.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the Target Group with Id: %s, Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.TargetGroups().Delete(*id)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func getTargetGroupNew(c *core.CommandConfig) resources.TargetGroup {
	input := resources.TargetGroupProperties{}
	// Set Required Properties
	input.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))))

	input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm))))

	input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))

	inputHealthCheck := resources.TargetGroupHealthCheck{}

	// Set Properties for Health Check for Target Group
	inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout))))

	inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))))

	inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries))))

	// Set Health Check for Target Group
	input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting HealthCheck"))

	inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}
	// Set Properties for Http Health Check for Target Group
	inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath))))

	inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))))

	inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType))))

	inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse))))

	inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex))))

	inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate))))

	// Set Http Health Check for Target Group
	input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Setting HttpHealthCheck"))

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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)) {
		input.SetAlgorithm(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Algorithm set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAlgorithm))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)) {
		input.SetProtocol(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property Protocol set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgProtocol))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
		inputHealthCheck := resources.TargetGroupHealthCheck{}

		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)) {
			inputHealthCheck.SetCheckTimeout(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout)))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property CheckTimeout for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckTimeout))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)) {
			inputHealthCheck.SetCheckInterval(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property CheckInterval for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgCheckInterval))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)) {
			inputHealthCheck.SetRetries(viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Retries for HealthCheck set: %v", viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgRetries))))
		}

		input.SetHealthCheck(inputHealthCheck.TargetGroupHealthCheck)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating HealthCheck"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
		inputHttpHealthCheck := resources.TargetGroupHttpHealthCheck{}

		// Set new values for Health Check Properties
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPath)) {
			inputHttpHealthCheck.SetPath(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Path for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPath))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) {
			inputHttpHealthCheck.SetMethod(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Method for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)) {
			inputHttpHealthCheck.SetResponse(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Response for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResponse))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)) {
			inputHttpHealthCheck.SetMatchType(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property MatchType for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMatchType))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)) {
			inputHttpHealthCheck.SetRegex(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRegex)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Regex for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRegex))))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)) {
			inputHttpHealthCheck.SetNegate(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNegate)))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property Negate for HttpHealthCheck set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNegate))))
		}

		input.SetHttpHealthCheck(inputHttpHealthCheck.TargetGroupHttpHealthCheck)
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Updating HttpHealthCheck"))
	}

	return &input
}
