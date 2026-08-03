package user

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func UserCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "user",
		Resource:  "user",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a new User on the contract",
		LongDesc: `Create a new User on your contract. First name, last name, email and password are required; the email address is the User's login and must be unique across IONOS CLOUD.

A newly created User has NO permissions by default. Give them access either by making them an administrator (--administrator, full contract access) or - the usual approach - by adding them to one or more Groups afterwards with ` + "`ionosctl compute group user add`" + `, so they inherit those groups' privileges.

Required values to run a command:

* First Name
* Last Name
* Email
* Password`,
		Example: `# Create a standard user (grant permissions later via group membership)
ionosctl compute user create --first-name Jane --last-name Doe --email jane.doe@example.com --password 's3cr3tPw'

# Create a full contract administrator who also must use two-factor auth
ionosctl compute user create --first-name Admin --last-name User --email admin@example.com --password 's3cr3tPw' --administrator --force-secure-auth`,
		PreCmdRun:  PreRunUserNameEmailPwd,
		CmdRun:     RunUserCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgFirstName, "", "", "The User's first name", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgLastName, "", "", "The User's last name", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The User's email address. This is the login identity and must be unique across IONOS CLOUD", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "The User's initial password (at least 5 characters). The User can change it after first login", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Make the User a contract administrator with full access to everything on the contract, bypassing all group privileges. Leave false for a normal User whose access comes from group membership")
	cmd.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Force two-factor (secure) authentication for this User: they must set up 2FA before they can sign in")

	return cmd
}
