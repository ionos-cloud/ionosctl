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
			Use:     "user",
			Short:   "Mongo Users Operations",
			Aliases: []string{"u"},
			Long: `The sub-commands of ` + "`ionosctl dbaas mongo user`" + ` manage the database users of a MongoDB cluster (separate from the IONOS Cloud account that owns the cluster).

Each user has a username, a password, and a set of ROLES. A role is a (database, role-name) pair that grants a privilege on a specific database - MongoDB authorization is per-database. Built-in role names include:
  read, readWrite                         - read-only / read+write on one database
  dbAdmin                                 - schema/index/stats admin on one database
  readAnyDatabase, readWriteAnyDatabase   - the same, across ALL databases
  dbAdminAnyDatabase                      - dbAdmin across ALL databases
  clusterMonitor                          - read-only monitoring of the whole cluster
  enableSharding                          - allow sharding operations on a database

A user is scoped to the cluster it is created in (--cluster-id).`,
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
