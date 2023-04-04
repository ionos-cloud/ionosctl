package user

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
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
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
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

type UserPrint struct {
	Username  string `json:"Username,omitempty"`
	CreatedBy string `json:"CreatedBy,omitempty"`
	Roles     string `json:"Roles,omitempty"`
}

var allCols = structs.Names(UserPrint{})

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

func getUserRows(ls *[]ionoscloud.User) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*ls))
	for _, t := range *ls {
		var cols UserPrint
		properties, ok := t.GetPropertiesOk()
		if ok {
			rolesAsStrings := functional.Map(*properties.GetRoles(), roleToString)
			cols.Roles = strings.Join(rolesAsStrings, ", ") // "db1: read, db2: write, db3: abcd..."

			cols.Username = *properties.GetUsername()
		}
		metadata, ok := t.GetMetadataOk()
		if ok {
			cols.CreatedBy = *metadata.GetCreatedBy()
		}
		o := structs.Map(cols)
		out = append(out, o)
	}
	return out
}

func getUserPrint(c *core.CommandConfig, ls *[]ionoscloud.User) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil && ls != nil {
		r.OutputJSON = ls
		r.KeyValue = getUserRows(ls)                            // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
	}
	return r
}
