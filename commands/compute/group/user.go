package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
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
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "Group User Operations",
			Long:             "The sub-commands of `ionosctl compute group user` allow you to list, add, remove Users from a Group.",
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
		ShortDesc:  "List Users from a Group",
		LongDesc:   "Use this command to get a list of Users from a specific Group.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.UsersFiltersUsage(),
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
		ShortDesc:  "Add User to a Group",
		LongDesc:   "Use this command to add an existing User to a specific Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
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
		ShortDesc:  "Remove User from a Group",
		LongDesc:   "Use this command to remove a User from a Group.\n\nRequired values to run command:\n\n* Group Id\n* User Id",
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
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Remove all Users from a group.")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for User removal to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for User removal [seconds]")

	return cmd
}
