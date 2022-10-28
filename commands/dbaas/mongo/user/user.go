package user

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
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
	return cmd
}

type UserPrint struct {
	UserId      string `json:"UserId,omitempty"`
	Cores       int32  `json:"Cores,omitempty"`
	StorageSize string `json:"StorageSize,omitempty"`
	Ram         string `json:"Ram,omitempty"`
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

func getUserRows(ls *[]ionoscloud.User) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*ls))
	for _, t := range *ls {
		var cols UserPrint
		properties, ok := t.GetPropertiesOk()
		if !ok {
			continue
		}
		roles := properties.GetRoles()

		rolesAsStrings :=
		for _, r := range *roles {
			roleToString(r)
		}

		o := structs.Map(cols)
		out = append(out, o)
	}
	return out
}

func getClusterHeaders(customColumns []string) []string {
	if customColumns == nil {
		return allCols[0:6]
	}
	//for _, c := customColumns {
	//	if slices.Contains(allCols, c) {}
	//}
	return customColumns
}
