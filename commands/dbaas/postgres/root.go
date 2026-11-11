package postgres

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/backup"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/version"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func Root() *core.Command {
	return newRoot("postgres", []string{"pg", "pgsql", "postgresql", "psql"})
}

// HiddenV2Root keeps the pre-promotion `postgres-v2` tree reachable so existing scripts keep working.
func HiddenV2Root() *core.Command {
	cmd := newRoot("postgres-v2", []string{"pg-v2", "pgsql-v2", "postgresql-v2", "psql-v2"})
	cmd.Command.Hidden = true
	return cmd
}

func newRoot(use string, aliases []string) *core.Command {
	pgsqlCmd := &core.Command{
		Command: &cobra.Command{
			Use:     use,
			Aliases: aliases,
			Short:   "Manage DBaaS PostgreSQL clusters and backups",
			Long: `Manage IONOS DBaaS PostgreSQL, a fully managed PostgreSQL service.

A ` + "`cluster`" + ` is the top-level resource: one or more PostgreSQL instances provisioned in a physical location, reachable over a private LAN in one of your virtual datacenters. Clusters support automated point-in-time backups used for restore and cloning.

PostgreSQL is a regional service, so commands operate against a specific location; set it with --location (or the config) where required.`,
			TraverseChildren: true,
		},
	}
	pgsqlCmd.AddCommand(cluster.ClusterCmd())
	pgsqlCmd.AddCommand(version.VersionCmd())
	pgsqlCmd.AddCommand(backup.BackupCmd())

	// Config key stays `psqlv2`: the v1 `psql` key points at the decommissioned v1 API.
	return core.WithRegionalConfigOverride(pgsqlCmd, []string{fileconfiguration.PSQLV2}, constants.PostgresApiRegionalURL, constants.PostgresLocations)
}
