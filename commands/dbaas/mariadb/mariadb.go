package mariadb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	mongoCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "mariadb",
			Aliases:          []string{"maria", "mar"},
			Short:            "DBaaS MariaDB Operations",
			TraverseChildren: true,
		},
	}
	mongoCmd.AddCommand(cluster.Root())
	mongoCmd.AddCommand(backup.Root())
	return mongoCmd
}
