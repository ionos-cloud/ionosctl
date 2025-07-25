package mariadb

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "mariadb",
			Aliases:          []string{"maria", "mar", "ma"},
			Short:            "DBaaS MariaDB Operations",
			TraverseChildren: true,
		},
	}
	cmd.AddCommand(cluster.Root())
	cmd.AddCommand(backup.Root())

	return core.WithRegionalConfigOverride(cmd, []string{"mariadb"}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations)
}
