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
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "User Operations",
			Long:             "The sub-commands of `ionosctl compute user` allow you to list, get, create, update, delete Users under your account. To add Users to a Group, check the `ionosctl compute group user` commands. To add S3Keys to a User, check the `ionosctl compute user s3key` commands.",
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
