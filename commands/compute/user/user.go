package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
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

func UserCmd() *core.Command {
	userCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "user",
			Aliases: []string{"u"},
			Short:   "Manage the cloud user accounts on your contract",
			Long: `A User is an individual cloud account on your contract - identified by an email address, holding a password, and able to sign in to the DCD web UI and the API. Users are the "who" of IONOS Cloud Identity & Access Management (IAM).

A User gets permissions in one of two ways:
  * Administrator: setting --administrator makes the User a full contract admin who bypasses all group privileges and can do anything on the contract. Use sparingly.
  * Group membership: a non-admin User's permissions are the UNION of the privileges of every Group they belong to. Add Users to Groups with ` + "`ionosctl compute group user add`" + `.

Related sub-trees: ` + "`ionosctl compute group user`" + ` manages which Groups a User belongs to; ` + "`ionosctl compute user s3key`" + ` manages a User's Object-Storage (S3-compatible) access keys.`,
			TraverseChildren: true,
		},
	}

	userCmd.AddColsFlag(allUserCols)

	userCmd.AddCommand(UserListCmd())
	userCmd.AddCommand(UserGetCmd())
	userCmd.AddCommand(UserCreateCmd())
	userCmd.AddCommand(UserUpdateCmd())
	userCmd.AddCommand(UserDeleteCmd())
	userCmd.AddCommand(S3keyCmd())

	return core.WithConfigOverride(userCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
