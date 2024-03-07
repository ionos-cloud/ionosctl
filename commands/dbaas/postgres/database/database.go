package database

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Name", "Owner", "ClusterId"}
	defaultCols = []string{"Id", "Name", "Owner"}
)

func DatabaseCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "database",
			Aliases:          []string{"databases"},
			Short:            "DBaaS Postgresql Database Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres database` allow you to perform operations on DBaaS PostgreSQL databases.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(ListCmd())
	cmd.AddCommand(GetCmd())
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(DeleteCmd())
	return cmd
}
