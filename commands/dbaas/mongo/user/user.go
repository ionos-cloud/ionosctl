package user

import (
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

const (
	FlagDatabase      = "database"
	FlagDatabaseShort = "d"
	FlagRoles         = "roles"
	FlagRolesShort    = "r"
)

func UserCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Short:            "Mongo Users Operations",
			Aliases:          []string{"u"},
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(UserListCmd())
	cmd.AddCommand(UserCreateCmd())
	cmd.AddCommand(UserGetCmd())
	cmd.AddCommand(UserDeleteCmd())
	return cmd
}

var (
	allCols = []string{"Username", "CreatedBy", "Roles"}
)
