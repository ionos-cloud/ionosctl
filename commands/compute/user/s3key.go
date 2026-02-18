package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultS3KeyCols = []string{"S3KeyId", "Active", "SecretKey"}
)

func S3keyCmd() *core.Command {
	s3keyCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "s3key",
			Aliases:          []string{"k", "s3k"},
			Short:            "User S3Key Operations",
			Long:             "The sub-commands of `ionosctl user s3key` allow you to see information, to list, get, create, update, delete Users S3Keys.",
			TraverseChildren: true,
		},
	}
	globalFlags := s3keyCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultS3KeyCols, tabheaders.ColsMessage(defaultS3KeyCols))
	_ = viper.BindPFlag(core.GetFlagName(s3keyCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = s3keyCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultS3KeyCols, cobra.ShellCompDirectiveNoFileComp
	})

	s3keyCmd.AddCommand(s3keyListCmd())
	s3keyCmd.AddCommand(s3keyGetCmd())
	s3keyCmd.AddCommand(s3keyCreateCmd())
	s3keyCmd.AddCommand(s3keyUpdateCmd())
	s3keyCmd.AddCommand(s3keyDeleteCmd())

	return core.WithConfigOverride(s3keyCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func s3keyListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List User S3Keys",
		LongDesc:   "Use this command to get a list of S3Keys of a specified User.\n\nRequired values to run command:\n\n* User Id",
		Example:    "ionosctl user s3key list --user-id USER_ID",
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserS3KeyList,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func s3keyGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a User S3Key",
		LongDesc:   "Use this command to get information about a specified S3Key from a specified User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    "ionosctl user s3key get --user-id USER_ID --s3key-id S3KEY_ID",
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func s3keyCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "user",
		Resource:  "s3key",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a S3Key for a User",
		LongDesc: `Use this command to create a S3Key for a particular User.

Note: A maximum of five S3 keys may be created for any given user.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* User Id`,
		Example:    "ionosctl user s3key create --user-id USER_ID",
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserS3KeyCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for User S3Key creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key creation [seconds]")

	return cmd
}

func s3keyUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "user",
		Resource:  "s3key",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a S3Key",
		LongDesc: `Use this command to update a specified S3Key from a particular User. This operation allows you to enable or disable a specific S3Key.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* User Id
* S3Key Id
* S3Key Active`,
		Example:    "ionosctl user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false",
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgS3KeyActive, "", false, "Enable or disable an User S3Key. E.g.: --s3key-active=true, --s3key-active=false")
	cmd.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for User S3Key update to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key update [seconds]")

	return cmd
}

func s3keyDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete a S3Key",
		LongDesc:   "Use this command to delete a specific S3Key of an User.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    "ionosctl user s3key delete --user-id USER_ID --s3key-id S3KEY_ID --force",
		PreCmdRun:  PreRunUserKeyDelete,
		CmdRun:     RunUserS3KeyDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for Request for User S3Key deletion to be executed")
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the S3Keys of an User.")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for User S3Key deletion [seconds]")

	return cmd
}
