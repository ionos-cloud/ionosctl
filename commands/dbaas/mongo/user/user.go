package user

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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

	cmd.AddColsFlag(allCols)

	cmd.AddCommand(UserListCmd())
	cmd.AddCommand(UserCreateCmd())
	cmd.AddCommand(UserGetCmd())
	cmd.AddCommand(UserDeleteCmd())
	return cmd
}

var allCols = []table.Column{
	{Name: "Username", JSONPath: "properties.username", Default: true},
	{Name: "CreatedBy", JSONPath: "metadata.createdBy", Default: true},
	{Name: "Roles", Default: true, Format: func(item map[string]any) any {
		roles, ok := table.Navigate(item, "properties.roles").([]any)
		if !ok || len(roles) == 0 {
			return nil
		}
		var parts []string
		for _, r := range roles {
			role, ok := r.(map[string]any)
			if !ok {
				continue
			}
			db, _ := role["database"].(string)
			roleName, _ := role["role"].(string)
			parts = append(parts, fmt.Sprintf("%s: %s", db, roleName))
		}
		return strings.Join(parts, ", ")
	}},
}
