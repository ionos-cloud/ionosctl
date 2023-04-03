package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultGroupShareCols, printer.ColsMessage(allGroupShareCols))
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
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
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
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
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
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ResourcesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
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
		return completer.GroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", cloudapiv6.ResourceId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Resource Share deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the Resources Share from a specified Group.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return shareCmd
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
	var allShares []resources.GroupShare
	totalTime := time.Duration(0)
	for _, group := range getGroups(groups) {
		shares, resp, err := c.CloudApiV6Services.Groups().ListShares(*group.GetId(), listQueryParams)
		if err != nil {
			return err
		}
		allShares = append(allShares, getGroupShares(shares)...)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Printer.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer.Print(getGroupSharePrint(nil, c, allShares))
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
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShares(shares)))
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
	c.Printer.Verbose("Getting Share with Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	s, resp, err := c.CloudApiV6Services.Groups().GetShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShare(s)))
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
	c.Printer.Verbose("Properties set for creating the Share: EditPrivilege: %v, SharePrivilege: %v", editPrivilege, sharePrivilege)
	c.Printer.Verbose("Adding Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	shareAdded, resp, err := c.CloudApiV6Services.Groups().AddShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareAdded)))
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
	c.Printer.Verbose("Updating Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	shareUpdated, resp, err := c.CloudApiV6Services.Groups().UpdateShare(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		newShare,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareUpdated)))
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
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete share from group"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Share with Resource ID: %v from Group with ID: %v...", shareId, groupId)
		resp, err := c.CloudApiV6Services.Groups().RemoveShare(groupId, shareId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getGroupSharePrint(resp, c, nil))
	}
}

func getShareUpdateInfo(oldShare *resources.GroupShare, c *core.CommandConfig) *resources.GroupShareProperties {
	var sharePrivilege, editPrivilege bool
	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgEditPrivilege))
			c.Printer.Verbose("Property EditPrivilege set: %v", editPrivilege)
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgSharePrivilege))
			c.Printer.Verbose("Property SharePrivilege set: %v", sharePrivilege)
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
	c.Printer.Verbose("Group ID: %v", groupId)
	c.Printer.Verbose("Getting Group Shares...")
	groupShares, resp, err := c.CloudApiV6Services.Groups().ListShares(groupId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if groupSharesItems, ok := groupShares.GetItemsOk(); ok && groupSharesItems != nil {
		if len(*groupSharesItems) > 0 {
			_ = c.Printer.Warn("GroupShares to be deleted:")
			for _, share := range *groupSharesItems {
				if id, ok := share.GetIdOk(); ok && id != nil {
					_ = c.Printer.Warn("GroupShare Id: " + *id)
				}
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the GroupShares"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the GroupShares...")
			var multiErr error
			for _, share := range *groupSharesItems {
				if id, ok := share.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Share with Resource ID: %v from Group with ID: %v...",
						*id, groupId)
					resp, err = c.CloudApiV6Services.Groups().RemoveShare(groupId, *id, queryParams)
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
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Group Shares found")
		}
	} else {
		return errors.New("could not get items of Group Shares")
	}
}

// Output Printing

var defaultGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type"}
var allGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type", "GroupId"}

type groupSharePrint struct {
	ShareId        string `json:"ShareId,omitempty"`
	EditPrivilege  bool   `json:"EditPrivilege,omitempty"`
	SharePrivilege bool   `json:"SharePrivilege,omitempty"`
	Type           string `json:"Type,omitempty"`
	GroupId        string `json:"GroupId,omitempty"`
}

func getGroupSharePrint(resp *resources.Response, c *core.CommandConfig, groups []resources.GroupShare) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
		}
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getGroupSharesKVMaps(groups)
			r.Columns = printer.GetHeadersListAll(allGroupShareCols, defaultGroupShareCols, "GroupId", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgCols)), viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)))
		}
	}
	return r
}

func getGroupShares(groups resources.GroupShares) []resources.GroupShare {
	u := make([]resources.GroupShare, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.GroupShare{GroupShare: item})
		}
	}
	return u
}

func getGroupShare(u *resources.GroupShare) []resources.GroupShare {
	groups := make([]resources.GroupShare, 0)
	if u != nil {
		groups = append(groups, resources.GroupShare{GroupShare: u.GroupShare})
	}
	return groups
}

func getGroupSharesKVMaps(gs []resources.GroupShare) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(gs))
	for _, g := range gs {
		var gPrint groupSharePrint
		if id, ok := g.GetIdOk(); ok && id != nil {
			gPrint.ShareId = *id
		}
		if properties, ok := g.GetPropertiesOk(); ok && properties != nil {
			if edit, ok := properties.GetEditPrivilegeOk(); ok && edit != nil {
				gPrint.EditPrivilege = *edit
			}
			if sh, ok := properties.GetSharePrivilegeOk(); ok && sh != nil {
				gPrint.SharePrivilege = *sh
			}
		}
		if typeResource, ok := g.GetTypeOk(); ok && typeResource != nil {
			gPrint.Type = string(*typeResource)
		}
		if hrefOk, ok := g.GetHrefOk(); ok && hrefOk != nil {
			// Get parent resource ID using HREF: `.../um/groups/[PARENT_ID_WE_WANT]/shares/[SHARE_ID]`
			gPrint.GroupId = strings.Split(strings.Split(*hrefOk, "groups")[1], "/")[1]
		}
		o := structs.Map(gPrint)
		out = append(out, o)
	}
	return out
}
