package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/backup/location"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
)

var backupCols = []table.Column{
	{Name: "BackupId", JSONPath: "id", Default: true},
	{Name: "ClusterId", JSONPath: "properties.clusterId", Default: true},
	{Name: "PostgresClusterVersion", JSONPath: "properties.postgresClusterVersion", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "IsActive", JSONPath: "properties.isActive", Default: true},
	{Name: "EarliestRecoveryTargetTime", JSONPath: "properties.earliestRecoveryTargetTime", Default: true},
	{Name: "LatestRecoveryTargetTime", JSONPath: "properties.latestRecoveryTargetTime", Default: true},
	{Name: "State", JSONPath: "metadata.state"},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate"},
}

var allBackupCols = table.AllCols(backupCols)
var defaultBackupCols = table.DefaultCols(backupCols)

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
