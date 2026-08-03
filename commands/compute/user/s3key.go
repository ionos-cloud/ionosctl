package user

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

var allS3KeyCols = []table.Column{
	{Name: "S3KeyId", JSONPath: "id", Default: true},
	{Name: "Active", JSONPath: "properties.active", Default: true},
	{Name: "SecretKey", JSONPath: "properties.secretKey", Default: true},
}

func S3keyCmd() *core.Command {
	s3keyCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "s3key",
			Aliases: []string{"k", "s3k"},
			Short:   "Manage a User's Object-Storage (S3) access keys",
			Long: `Manage the S3-compatible access keys belonging to a User. Each key is an access-key-ID / secret-key pair that programs (aws-cli, s3cmd, SDKs, etc.) use to authenticate against IONOS Object Storage on that User's behalf.

The secret key is only visible via get/list; treat it like a password. Keys can be individually disabled (--s3key-active=false) without deleting them, which is the recommended way to rotate credentials. Using Object Storage at all also requires the User (or one of their Groups) to hold the S3 privilege (see --s3privilege on ` + "`ionosctl compute group`" + `).

Note: a maximum of five S3 keys may exist per User at any time.`,
			TraverseChildren: true,
		},
	}
	s3keyCmd.AddColsFlag(allS3KeyCols)

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
		ShortDesc:  "List a User's S3 keys",
		LongDesc:   "List every S3 access key belonging to the given User, including each key's active/disabled state and its secret key.\n\nRequired values to run command:\n\n* User Id",
		Example:    "ionosctl compute user s3key list --user-id USER_ID",
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
		ShortDesc:  "Get one of a User's S3 keys",
		LongDesc:   "Retrieve a single S3 key of a User, including its secret key (needed to configure an S3 client). The key ID is the access-key-ID.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    "ionosctl compute user s3key get --user-id USER_ID --s3key-id S3KEY_ID",
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
		ShortDesc: "Generate a new S3 key for a User",
		LongDesc: `Generate a new S3 access-key pair for the given User. The API returns both the access-key-ID and the secret key; the secret is retrievable later via get/list, but you should still capture it now and store it securely. The key is created in the active (enabled) state.

Note: a maximum of five S3 keys may exist for any given User. If the User already has five, delete or reuse one before creating another.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* User Id`,
		Example:    "ionosctl compute user s3key create --user-id USER_ID",
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserS3KeyCreate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func s3keyUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "user",
		Resource:  "s3key",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Enable or disable a User's S3 key",
		LongDesc: `Enable or disable an existing S3 key of a User by setting --s3key-active. Disabling a key immediately stops it from authenticating against Object Storage without deleting it, which makes this the safe way to rotate or temporarily suspend credentials (re-enable it later, or delete it once a replacement is in use).

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* User Id
* S3Key Id
* S3Key Active`,
		Example:    "ionosctl compute user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false",
		PreCmdRun:  PreRunUserKeyIds,
		CmdRun:     RunUserS3KeyUpdate,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, "", "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgS3KeyActive, "", false, "Whether the key is active: true enables it for Object-Storage authentication, false disables it (without deleting it). E.g.: --s3key-active=true, --s3key-active=false")
	cmd.AddStringFlag(cloudapiv6.ArgS3KeyId, cloudapiv6.ArgIdShort, "", cloudapiv6.S3KeyId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgS3KeyId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.S3KeyIds(viper.GetString(core.GetFlagName(cmd.NS, cloudapiv6.ArgUserId))), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}

func s3keyDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "s3key",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Permanently delete a User's S3 key",
		LongDesc:   "Permanently delete a specific S3 key of a User. Any client still configured with this key immediately loses access to Object Storage; there is no undo. To only pause a key, disable it with `s3key update --s3key-active=false` instead.\n\nRequired values to run command:\n\n* User Id\n* S3Key Id",
		Example:    "ionosctl compute user s3key delete --user-id USER_ID --s3key-id S3KEY_ID --force",
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
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every S3 key of the User, revoking all of their Object-Storage credentials at once")

	return cmd
}
