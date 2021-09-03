package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func share() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultGroupShareCols, utils.ColsMessage(defaultGroupShareCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(shareCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = shareCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultGroupShareCols, cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunShareList,
		InitClient: true,
	})
	list.AddStringFlag(config.ArgGroupId, "", "", config.GroupId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	get.AddStringFlag(config.ArgGroupId, "", "", config.GroupId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgResourceId, config.ArgIdShort, "", config.ResourceId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})

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
	create.AddStringFlag(config.ArgGroupId, "", "", config.GroupId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgResourceId, config.ArgIdShort, "", config.ResourceId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getResourcesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgEditPrivilege, "", false, "Set the group's permission to edit privileges on resource")
	create.AddBoolFlag(config.ArgSharePrivilege, "", false, "Set the group's permission to share resource")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Resource share to executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Resource to be shared through a Group [seconds]")

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
	update.AddStringFlag(config.ArgGroupId, "", "", config.GroupId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgResourceId, config.ArgIdShort, "", config.ResourceId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgEditPrivilege, "", false, "Update the group's permission to edit privileges on resource")
	update.AddBoolFlag(config.ArgSharePrivilege, "", false, "Update the group's permission to share resource")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Resource Share update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share update [seconds]")

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
		PreCmdRun:  PreRunGroupResourceIds,
		CmdRun:     RunShareDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgGroupId, "", "", config.GroupId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgResourceId, config.ArgIdShort, "", config.ResourceId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Resource Share deletion to be executed")
	deleteCmd.AddBoolFlag(config.ArgAll, config.ArgAllShort, false, "delete all the Resources Share from a specified Group.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share deletion [seconds]")

	return shareCmd
}

func PreRunGroupResourceIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgGroupId, config.ArgResourceId)
}

func RunShareList(c *core.CommandConfig) error {
	shares, resp, err := c.Groups().ListShares(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShares(shares)))
}

func RunShareGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Share with Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	s, resp, err := c.Groups().GetShare(
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShare(s)))
}

func RunShareCreate(c *core.CommandConfig) error {
	editPrivilege := viper.GetBool(core.GetFlagName(c.NS, config.ArgEditPrivilege))
	sharePrivilege := viper.GetBool(core.GetFlagName(c.NS, config.ArgSharePrivilege))
	input := v5.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &editPrivilege,
				SharePrivilege: &sharePrivilege,
			},
		},
	}
	c.Printer.Verbose("Properties set for creating the Share: EditPrivilege: %v, SharePrivilege: %v", editPrivilege, sharePrivilege)
	c.Printer.Verbose("Adding Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	shareAdded, resp, err := c.Groups().AddShare(
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareAdded)))
}

func RunShareUpdate(c *core.CommandConfig) error {
	s, _, err := c.Groups().GetShare(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)), viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)))
	if err != nil {
		return err
	}
	properties := getShareUpdateInfo(s, c)
	newShare := v5.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &properties.GroupShareProperties,
		},
	}
	c.Printer.Verbose("Updating Share for Resource ID: %v from Group with ID: %v...",
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	shareUpdated, resp, err := c.Groups().UpdateShare(
		viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		newShare,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareUpdated)))
}

func RunShareDelete(c *core.CommandConfig) error {
	var resp *v5.Response
	var err error
	var groupShares v5.GroupShares
	groupId := viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId))
	shareId := viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId))
	allFlag := viper.GetBool(core.GetFlagName(c.NS, config.ArgAll))
	if allFlag {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "Are you sure you want to delete all the GroupShares?"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the GroupShares...")
		groupShares, resp, err = c.Groups().ListShares(groupId)
		if err != nil {
			return err
		}
		if groupSharesItems, ok := groupShares.GetItemsOk(); ok && groupSharesItems != nil {
			for _, share := range *groupSharesItems {
				if id, ok := share.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Deleting GroupShare with id: %v...", *id)
					resp, err = c.Groups().RemoveShare(groupId, *id)
					if err != nil {
						return err
					}
					if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
						return err
					}
				}
			}
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete share from group"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting Share with Resource ID: %v from Group with ID: %v...",
			viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
			viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
		resp, err := c.Groups().RemoveShare(groupId, shareId)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
			return err
		}
	}

	return c.Printer.Print(getGroupSharePrint(resp, c, nil))
}

func getShareUpdateInfo(oldShare *v5.GroupShare, c *core.CommandConfig) *v5.GroupShareProperties {
	var sharePrivilege, editPrivilege bool
	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(core.GetFlagName(c.NS, config.ArgEditPrivilege))
			c.Printer.Verbose("Property EditPrivilege set: %v", editPrivilege)
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(core.GetFlagName(c.NS, config.ArgSharePrivilege))
			c.Printer.Verbose("Property SharePrivilege set: %v", sharePrivilege)
		} else {
			if e, ok := properties.GetSharePrivilegeOk(); ok && e != nil {
				sharePrivilege = *e
			}
		}
	}
	return &v5.GroupShareProperties{
		GroupShareProperties: ionoscloud.GroupShareProperties{
			EditPrivilege:  &editPrivilege,
			SharePrivilege: &sharePrivilege,
		},
	}
}

// Output Printing

var defaultGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type"}

type groupSharePrint struct {
	ShareId        string `json:"ShareId,omitempty"`
	EditPrivilege  bool   `json:"EditPrivilege,omitempty"`
	SharePrivilege bool   `json:"SharePrivilege,omitempty"`
	Type           string `json:"Type,omitempty"`
}

func getGroupSharePrint(resp *v5.Response, c *core.CommandConfig, groups []v5.GroupShare) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getGroupSharesKVMaps(groups)
			r.Columns = getGroupShareCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getGroupShareCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"ShareId":        "ShareId",
			"EditPrivilege":  "EditPrivilege",
			"SharePrivilege": "SharePrivilege",
			"Type":           "Type",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultGroupShareCols
	}
}

func getGroupShares(groups v5.GroupShares) []v5.GroupShare {
	u := make([]v5.GroupShare, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, v5.GroupShare{GroupShare: item})
		}
	}
	return u
}

func getGroupShare(u *v5.GroupShare) []v5.GroupShare {
	groups := make([]v5.GroupShare, 0)
	if u != nil {
		groups = append(groups, v5.GroupShare{GroupShare: u.GroupShare})
	}
	return groups
}

func getGroupSharesKVMaps(gs []v5.GroupShare) []map[string]interface{} {
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
		o := structs.Map(gPrint)
		out = append(out, o)
	}
	return out
}

func getGroupResourcesIds(outErr io.Writer, groupId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	groupSvc := v5.NewGroupService(clientSvc.Get(), context.TODO())
	res, _, err := groupSvc.ListResources(groupId)
	clierror.CheckError(err, outErr)
	resIds := make([]string, 0)
	if items, ok := res.ResourceGroups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				resIds = append(resIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return resIds
}
