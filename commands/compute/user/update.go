package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func UserUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a User's details or admin/2FA status",
		LongDesc: `Update a User's profile fields (name, email, password) and their contract-level flags (administrator, forced two-factor auth). Only the fields you pass are changed.

Note: this command does NOT change which Groups a User belongs to - manage membership (and therefore inherited group privileges) with ` + "`ionosctl compute group user add/remove`" + `. Toggling --administrator here is the one way to grant or revoke blanket contract-wide access directly on the User.

Required values to run command:

* User Id`,
		Example: `# Promote a user to contract administrator
ionosctl compute user update --user-id USER_ID --admin=true

# Demote an admin back to a normal user and require 2FA
ionosctl compute user update --user-id USER_ID --admin=false --force-secure-auth=true

# Reset a user's password
ionosctl compute user update --user-id USER_ID --password 'newS3cr3t'`,
		PreCmdRun:  PreRunUserId,
		CmdRun:     RunUserUpdate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgFirstName, "", "", "The User's first name")
	cmd.AddStringFlag(cloudapiv6.ArgLastName, "", "", "The User's last name")
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The User's email address (login identity). Must remain unique across IONOS Cloud")
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Set a new password for the User (at least 5 characters)")
	cmd.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Grant (true) or revoke (false) contract-administrator rights - full access to the whole contract, bypassing group privileges. E.g.: --admin=true, --admin=false")
	cmd.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Force (true) or stop forcing (false) two-factor authentication for this User. E.g.: --force-secure-auth=true, --force-secure-auth=false")
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
