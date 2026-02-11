package backup

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/backup/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultBackupCols = []string{"BackupId", "ClusterId", "PostgresClusterVersion", "Location", "IsActive", "EarliestRecoveryTargetTime", "LatestRecoveryTargetTime"}
	allBackupCols     = []string{"BackupId", "ClusterId", "PostgresClusterVersion", "Location", "IsActive", "EarliestRecoveryTargetTime", "LatestRecoveryTargetTime", "State", "CreatedDate"}
)

func BackupCmd() *core.Command {
	backupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "PostgreSQL Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres backup` allow you to list and get DBaaS PostgreSQL Backups.",
			TraverseChildren: true,
		},
	}

	backupCmd.AddCommand(BackupListCmd())
	backupCmd.AddCommand(BackupGetCmd())
	backupCmd.AddCommand(location.BackupLocationCmd())

	return backupCmd
}
