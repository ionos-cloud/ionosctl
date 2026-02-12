package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/backup/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
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

func Backups(fs ...Filter) (psqlv2.BackupReadList, error) {
	req := client.Must().PostgresClientV2.BackupsApi.BackupsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return psqlv2.BackupReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return psqlv2.BackupReadList{}, err
	}
	return ls, nil
}

type Filter func(request psqlv2.ApiBackupsGetRequest) (psqlv2.ApiBackupsGetRequest, error)
