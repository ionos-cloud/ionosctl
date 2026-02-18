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
		ShortDesc: "Create a User under a particular contract",
		LongDesc: `Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Required values to run a command:

* First Name
* Last Name
* Email
* Password`,
		Example:    "ionosctl user create --first-name NAME --last-name NAME --email EMAIL --password PASSWORD",
		PreCmdRun:  PreRunUserNameEmailPwd,
		CmdRun:     RunUserCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgFirstName, "", "", "The first name for the User", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgLastName, "", "", "The last name for the User", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The email for the User", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "The password for the User (must be at least 5 characters long)", core.RequiredFlagOption())
	cmd.AddBoolFlag(cloudapiv6.ArgAdmin, "", false, "Assigns the User to have administrative rights")
	cmd.AddBoolFlag(cloudapiv6.ArgForceSecAuth, "", false, "Indicates if secure (two-factor) authentication should be forced for the User")

	return cmd
}
