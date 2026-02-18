package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultUserCols = []string{"UserId", "Firstname", "Lastname", "Email", "S3CanonicalUserId", "Administrator", "ForceSecAuth", "SecAuthActive", "Active"}
)

func UserCmd() *core.Command {
	userCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "User Operations",
			Long:             "The sub-commands of `ionosctl user` allow you to list, get, create, update, delete Users under your account. To add Users to a Group, check the `ionosctl group user` commands. To add S3Keys to a User, check the `ionosctl user s3key` commands.",
			TraverseChildren: true,
		},
	}

	globalFlags := userCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultUserCols, tabheaders.ColsMessage(defaultUserCols))
	_ = userCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultUserCols, cobra.ShellCompDirectiveNoFileComp
	})

	userCmd.AddCommand(UserListCmd())
	userCmd.AddCommand(UserGetCmd())
	userCmd.AddCommand(UserCreateCmd())
	userCmd.AddCommand(UserUpdateCmd())
	userCmd.AddCommand(UserDeleteCmd())
	userCmd.AddCommand(S3keyCmd())

	return core.WithConfigOverride(userCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
