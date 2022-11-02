package user

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDatabase = "database"
)

func UserCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"u"},
			Short:            "Mongo Users Operations",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(UserListCmd())
	cmd.AddCommand(UserGetCmd())
	cmd.AddCommand(UserDeleteCmd())
	return cmd
}

type UserPrint struct {
	Username  string `json:"Username,omitempty"`
	Database  string `json:"Database,omitempty"`
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
			rolesAsStrings := utils.MapNoIdx(*properties.GetRoles(), roleToString)
			cols.Roles = strings.Join(rolesAsStrings, ", ") // "db1: read, db2: write, db3: abcd..."

			cols.Database = *properties.GetDatabase()
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
	if c != nil && ls != nil {
		r.OutputJSON = ls
		r.KeyValue = getUserRows(ls)                                                                                       // map header -> rows
		r.Columns = printer.GetHeadersAllDefault(allCols, viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
	}
	return r
}
