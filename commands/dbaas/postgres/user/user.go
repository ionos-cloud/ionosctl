package user

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Username", "System", "ClusterId"}
	defaultCols = []string{"Id", "Username", "System"}
)

func UserCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "user",
			Aliases:          []string{"usr", "u", "users"},
			Short:            "DBaaS Postgresql User Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres user` allow you to perform operations on DBaaS PostgreSQL users.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(ListCmd())
	return cmd
}
