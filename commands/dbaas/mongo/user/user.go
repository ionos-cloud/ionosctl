package user

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
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
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	cmd.AddCommand(UserListCmd())
	cmd.AddCommand(UserCreateCmd())
	cmd.AddCommand(UserGetCmd())
	cmd.AddCommand(UserDeleteCmd())
	return cmd
}

var (
	allCols = []string{"Username", "CreatedBy", "Roles"}
)

func convertUserToTable(user ionoscloud.User) ([]map[string]interface{}, error) {
	properties, ok := user.GetPropertiesOk()
	if !ok || properties == nil {
		return nil, fmt.Errorf("could not retrieve Mongo User properties")
	}

	roles, ok := properties.GetRolesOk()
	if !ok || roles == nil {
		return nil, fmt.Errorf("could not retrieve Mongo User roles")
	}

	temp, err := json2table.ConvertJSONToTable("", jsonpaths.User, user)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["Roles"] = strings.Join(functional.Map(*properties.GetRoles(), roleToString), ", ")

	return temp, nil
}

func convertUsersToTable(users ionoscloud.UsersList) ([]map[string]interface{}, error) {
	items, ok := users.GetItemsOk()
	if !ok || items == nil {
		return nil, fmt.Errorf("could not retrieve Mongo Users items")
	}

	var usersConverted []map[string]interface{}
	for _, item := range *items {
		temp, err := convertUserToTable(item)
		if err != nil {
			return nil, err
		}

		usersConverted = append(usersConverted, temp...)
	}

	return usersConverted, nil
}

// given a User DB/Role pair, return its string representation
// Role: { "role": "read", "database": "db" } -> "db: read"
func roleToString(role ionoscloud.UserRoles) string {
	val, ok := role.GetRoleOk()
	if !ok {
		return ""
	}
	db, ok := role.GetDatabaseOk()
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s: %s", *db, *val)
}
