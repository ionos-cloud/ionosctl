package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allUserCols = []table.Column{
	{Name: "UserId", JSONPath: "id", Default: true},
	{Name: "Firstname", JSONPath: "properties.firstName", Default: true},
	{Name: "Lastname", JSONPath: "properties.lastName", Default: true},
	{Name: "Email", JSONPath: "properties.email", Default: true},
	{Name: "S3CanonicalUserId", JSONPath: "properties.s3CanonicalUserId", Default: true},
	{Name: "Administrator", JSONPath: "properties.administrator", Default: true},
	{Name: "ForceSecAuth", JSONPath: "properties.forceSecAuth", Default: true},
	{Name: "SecAuthActive", JSONPath: "properties.secAuthActive", Default: true},
	{Name: "Active", JSONPath: "properties.active", Default: true},
}

func GroupUserCmd() *core.Command {
	groupUserCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "user",
			Aliases: []string{"u"},
			Short:   "Manage Group membership (which Users belong to a Group)",
			Long: `Manage the membership of an IAM Group: which Users belong to it. Membership is the link between Users and privileges - the moment a User is added to a Group, they inherit all of that Group's privileges (the union across every Group they belong to); removing them takes those privileges away (unless another Group still grants them).

These commands only add or remove existing Users to/from a Group. Create the Users first with ` + "`ionosctl compute user create`" + `, and set the Group's privileges with ` + "`ionosctl compute group create/update`" + `.`,
			TraverseChildren: true,
		},
	}

	groupUserCmd.AddCommand(groupUserListCmd())
	groupUserCmd.AddCommand(groupUserAddCmd())
	groupUserCmd.AddCommand(groupUserRemoveCmd())

	return core.WithConfigOverride(groupUserCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func groupUserListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List the Users that belong to a Group",
		LongDesc:   "List every User who is a member of the given Group (and therefore inherits its privileges).\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.UsersFiltersUsage(),
		Example:    "ionosctl compute group user list --group-id GROUP_ID",
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUserList,
		InitClient: true,
	})
	cmd.AddColsFlag(allUserCols)
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func groupUserAddCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "add",
		Aliases:    []string{"a"},
		ShortDesc:  "Add an existing User to a Group",
		LongDesc:   "Add an existing User to a Group. On success the User immediately inherits all of the Group's privileges and gains access to every resource shared with the Group. The User must already exist (create one with `ionosctl compute user create`).\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    "ionosctl compute group user add --group-id GROUP_ID --user-id USER_ID",
		PreCmdRun:  PreRunGroupUserIds,
		CmdRun:     RunGroupUserAdd,
		InitClient: true,
	})
	cmd.AddColsFlag(allUserCols)
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func groupUserRemoveCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "user",
		Verb:       "remove",
		Aliases:    []string{"r"},
		ShortDesc:  "Remove a User from a Group",
		LongDesc:   "Remove a User from a Group. The User keeps existing, but loses every privilege and shared-resource access that this Group granted (unless another Group they belong to still grants them).\n\nRequired values to run command:\n\n* Group Id\n* User Id",
		Example:    "ionosctl compute group user remove --group-id GROUP_ID --user-id USER_ID",
		PreCmdRun:  PreRunGroupUserRemove,
		CmdRun:     RunGroupUserRemove,
		InitClient: true,
	})
	cmd.AddColsFlag(allUserCols)
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupUsersIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgGroupId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove every User from the Group, leaving it empty. The Users and its privileges are not deleted")

	return cmd
}
