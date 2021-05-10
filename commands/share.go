package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func share() *builder.Command {
	ctx := context.TODO()
	shareCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "share",
			Short:            "Resource Share Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl share` + "`" + ` allow you to list, get, create, update, delete Resource Shares.`,
			TraverseChildren: true,
		},
	}
	globalFlags := shareCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultGroupShareCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(shareCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	list := builder.NewCommand(ctx, shareCmd, PreRunGroupId, RunShareList, "list", "List Resources Shares through a Group",
		"Use this command to get a full list of all the Resources that are shared through a specified Group.\n\nRequired values to run command:\n\n* Group Id",
		listSharesExample, true)
	list.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, shareCmd, PreRunGroupResourceIds, RunShareGet, "get", "Get a Resource Share from a Group",
		"Use this command to retrieve the details of a specific Shared Resource available to a specified Group.\n\nRequired values to run command:\n\n* Group Id\n* Resource Id",
		getShareExample, true)
	get.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgResourceId, "", "", config.RequiredFlagResourceId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(builder.GetFlagName(shareCmd.Name(), get.Name(), config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, shareCmd, PreRunGroupResourceIds, RunShareCreate, "create", "Create a Resource Share for a Group",
		`Use this command to create a specific Resource Share to a Group and optionally allow the setting of permissions for that Resource. As an example, you might use this to grant permissions to use an Image or Snapshot to a specific Group.

Required values to run a command:

* Group Id
* Resource Id`, createShareExample, true)
	create.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgResourceId, "", "", config.RequiredFlagResourceId)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getResourcesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgEditPrivilege, "", false, "Set the group's permission to edit privileges on resource")
	create.AddBoolFlag(config.ArgSharePrivilege, "", false, "Set the group's permission to share resource")
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Resource share to executed")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Resource to be shared through a Group [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, shareCmd, PreRunGroupResourceIds, RunShareUpdate, "update", "Update a Resource Share from a Group",
		`Use this command to update the permissions that a Group has for a specific Resource Share.

You can wait for the Request to be executed using `+"`"+`--wait-for-request`+"`"+` option.

Required values to run command:

* Group Id
* Resource Id`, updateShareExample, true)
	update.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgResourceId, "", "", config.RequiredFlagResourceId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(builder.GetFlagName(shareCmd.Name(), update.Name(), config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgEditPrivilege, "", false, "Update the group's permission to edit privileges on resource")
	update.AddBoolFlag(config.ArgSharePrivilege, "", false, "Update the group's permission to share resource")
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Resource Share update to be executed")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, shareCmd, PreRunGroupResourceIds, RunShareDelete, "delete", "Delete a Resource Share from a Group",
		`This command deletes a Resource Share from a specified Group.

Required values to run command:

* Resource Id
* Group Id`, deleteShareExample, true)
	deleteCmd.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgResourceId, "", "", config.RequiredFlagResourceId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupResourcesIds(os.Stderr, viper.GetString(builder.GetFlagName(shareCmd.Name(), deleteCmd.Name(), config.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for Resource Share deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Resource Share deletion [seconds]")

	return shareCmd
}

func PreRunGroupResourceIds(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgGroupId, config.ArgResourceId)
}

func RunShareList(c *builder.CommandConfig) error {
	shares, _, err := c.Groups().ListShares(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShares(shares)))
}

func RunShareGet(c *builder.CommandConfig) error {
	s, _, err := c.Groups().GetShare(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(nil, c, getGroupShare(s)))
}

func RunShareCreate(c *builder.CommandConfig) error {
	editPrivilege := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgEditPrivilege))
	sharePrivilege := viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgSharePrivilege))
	input := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &editPrivilege,
				SharePrivilege: &sharePrivilege,
			},
		},
	}
	shareAdded, resp, err := c.Groups().AddShare(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)),
		input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareAdded)))
}

func RunShareUpdate(c *builder.CommandConfig) error {
	s, _, err := c.Groups().GetShare(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)), viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)))
	if err != nil {
		return err
	}
	properties := getShareUpdateInfo(s, c)
	newShare := resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &properties.GroupShareProperties,
		},
	}
	shareUpdated, resp, err := c.Groups().UpdateShare(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)),
		newShare,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, getGroupShare(shareUpdated)))
}

func RunShareDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete share from group"); err != nil {
		return err
	}
	resp, err := c.Groups().RemoveShare(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)),
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupSharePrint(resp, c, nil))
}

func getShareUpdateInfo(oldShare *resources.GroupShare, c *builder.CommandConfig) *resources.GroupShareProperties {
	var sharePrivilege, editPrivilege bool
	if properties, ok := oldShare.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgEditPrivilege)) {
			editPrivilege = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgEditPrivilege))
		} else {
			if e, ok := properties.GetEditPrivilegeOk(); ok && e != nil {
				editPrivilege = *e
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgSharePrivilege)) {
			sharePrivilege = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgSharePrivilege))
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

// Output Printing

var defaultGroupShareCols = []string{"ShareId", "EditPrivilege", "SharePrivilege", "Type"}

type groupSharePrint struct {
	ShareId        string `json:"ShareId,omitempty"`
	EditPrivilege  bool   `json:"EditPrivilege,omitempty"`
	SharePrivilege bool   `json:"SharePrivilege,omitempty"`
	Type           string `json:"Type,omitempty"`
}

func getGroupSharePrint(resp *resources.Response, c *builder.CommandConfig, groups []resources.GroupShare) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitForRequest = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest))
		}
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getGroupSharesKVMaps(groups)
			r.Columns = getGroupShareCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
		o := structs.Map(gPrint)
		out = append(out, o)
	}
	return out
}

func getGroupResourcesIds(outErr io.Writer, groupId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(clientSvc.Get(), context.TODO())
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
