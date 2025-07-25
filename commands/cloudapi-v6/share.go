package commands

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type"}
	allGroupShareCols     = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type", "GroupId"}
)

func ShareCmd() *core.Command {
	ctx := context.TODO()
	shareCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "share",
			Short:            "Resource Share Operations",
			Long:             "The sub-commands of `ionosctl share` allow you to list, get, create, update, delete Resource Shares.",
			TraverseChildren: true,
		},
	}
	globalFlags := shareCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultGroupShareCols, tabheaders.ColsMessage(allGroupShareCols))
	_ = viper.BindPFlag(core.GetFlagName(shareCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = shareCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allGroupShareCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, shareCmd, core.CommandBuilder{
		Namespace:  "share",
		Resource:   "share",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Resources Shares through a Group",
		LongDesc:   "Use this command to get a full list of all the Resources that are shared through a specified Group.\n\nRequired values to run command:\n\n* Group Id",
		Example:    listSharesExample,
		PreCmdRun:  PreRunShareList,
		CmdRun:     RunShareList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, shareCmd, core.CommandBuilder{
		Namespace:  "share",
		Resource:   "share",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Resource Share from a Group",
		LongDesc:   "Use this command to retrieve the details of a specific Shared Resource available to a specified Group.\n\nRequired values to run command:\n\n* Group Id\n* Resource Id",
		Example:    getShareExample,
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, shareCmd, core.CommandBuilder{
		Namespace: "share",
		Resource:  "share",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Resource Share for a Group",
		LongDesc: `Use this command to create a specific Resource Share to a Group and optionally allow the setting of permissions for that Resource. As an example, you might use this to grant permissions to use an Image or Snapshot to a specific Group.

Required values to run a command:

* Group Id
* Resource Id`,
		Example:    createShareExample,
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareCreate,
		InitClient: true,
	})
	create.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ResourcesIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(cloudapiv6.ArgEditPrivilege, "", false, "Set the group's permission to edit privileges on resource")
	create.AddBoolFlag(cloudapiv6.ArgSharePrivilege, "", false, "Set the group's permission to share resource")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Resource share to executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Resource to be shared through a Group [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, shareCmd, core.CommandBuilder{
		Namespace: "share",
		Resource:  "share",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Resource Share from a Group",
		LongDesc: `Use this command to update the permissions that a Group has for a specific Resource Share.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Group Id
* Resource Id`,
		Example:    updateShareExample,
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(cloudapiv6.ArgEditPrivilege, "", false, "Update the group's permission to edit privileges on resource")
	update.AddBoolFlag(cloudapiv6.ArgSharePrivilege, "", false, "Update the group's permission to share resource")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Resource Share update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, shareCmd, core.CommandBuilder{
		Namespace: "share",
		Resource:  "share",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Resource Share from a Group",
		LongDesc: `This command deletes a Resource Share from a specified Group.

Required values to run command:

* Resource Id
* Group Id`,
		Example:    deleteShareExample,
		PreCmdRun:  PreRunGroupResourceDelete,
		CmdRun:     RunShareDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Resource Share deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Resources Share from a specified Group.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return core.WithConfigOverride(shareCmd, []string{"cloud", "compute"}, "")
}

func PreRunGroupResourceIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgGroupId, cloudapiv6.ArgResourceId)
}

func PreRunGroupResourceDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgResourceId},
		[]string{cloudapiv6.ArgGroupId, cloudapiv6.ArgAll},
	)
}

func RunShareListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	// Don't apply listQueryParams to parent resource, as it would have unexpected side effects on the results
	groups, _, err := c.CloudApiV6Services.Groups().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	var allShares = make([]ionoscloud.GroupShares, 0)

	totalTime := time.Duration(0)

	for _, group := range getGroups(groups) {
		shares, resp, err := c.CloudApiV6Services.Groups().ListShares(*group.GetId(), listQueryParams)
		if err != nil {
			return err
		}

		allShares = append(allShares, shares.GroupShares)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput(
		"*.items", jsonpaths.Share, allShares, tabheaders.GetHeaders(allGroupShareCols, defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunShareList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunShareListAll(c)
	}

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	shares, resp, err := c.CloudApiV6Services.Groups().ListShares(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Share, shares.GroupShares,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func PreRunShareList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ServersFilters(), completer.ServersFiltersUsage())
	}

	return nil
}

func RunShareGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Share with Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	s, resp, err := c.CloudApiV6Services.Groups().GetShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, s.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunShareCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	editPrivilege := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))
	sharePrivilege := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))

	input := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &editPrivilege,
				SharePrivilege: &sharePrivilege,
			},
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Share: EditPrivilege: %v, SharePrivilege: %v", editPrivilege, sharePrivilege))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Adding Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	shareAdded, resp, err := c.CloudApiV6Services.Groups().AddShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		input,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, shareAdded.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunShareUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	s, _, err := c.CloudApiV6Services.Groups().GetShare(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)), queryParams)
	if err != nil {
		return err
	}

	properties := getShareUpdateInfo(s, c)
	newShare := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &properties.GroupShareProperties,
		},
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Updating Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	shareUpdated, resp, err := c.CloudApiV6Services.Groups().UpdateShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		newShare,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Share, shareUpdated.GroupShare,
		tabheaders.GetHeadersAllDefault(defaultGroupShareCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunShareDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	shareId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId))
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllShares(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete share from group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Starting deleting Share with Resource ID: %v from Group with ID: %v...", shareId, groupId))

	resp, err := c.CloudApiV6Services.Groups().RemoveShare(groupId, shareId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Group Share successfully deleted"))
	return nil
}

func getShareUpdateInfo(oldShare *resources.GroupShare, c *core.CommandConfig) *resources.GroupShareProperties {
	var sharePrivilege, editPrivilege bool

	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property EditPrivilege set: %v", editPrivilege))
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SharePrivilege set: %v", sharePrivilege))
		} else {
			if e, ok := properties.GetSharePrivilegeOk(); ok && e != nil {
				sharePrivilege = *e
			}
		}
	}

	return &resources.GroupShareProperties{
		GroupShareProperties: ionoscloud.GroupShareProperties{
			EditPrivilege:  &editPrivilege,
			SharePrivilege: &sharePrivilege,
		},
	}
}

func DeleteAllShares(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Group ID: %v", groupId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Group Shares..."))

	groupShares, resp, err := c.CloudApiV6Services.Groups().ListShares(groupId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	groupSharesItems, ok := groupShares.GetItemsOk()
	if !ok || groupSharesItems == nil {
		return fmt.Errorf("could not get items of Group Shares")
	}

	if len(*groupSharesItems) <= 0 {
		return fmt.Errorf("no Group Shares found")
	}

	var multiErr error
	for _, share := range *groupSharesItems {
		id := share.GetId()
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the GroupShare with Id: %s", *id), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Groups().RemoveShare(groupId, *id, queryParams)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))

		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
