package backupunit

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

const backupUnitNote = "NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!"

var allCols = []table.Column{
	{Name: "BackupUnitId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Email", JSONPath: "properties.email", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

var allSSOUrlCols = []table.Column{
	{Name: "BackupUnitSsoUrl", JSONPath: "ssoUrl", Default: true},
}

func BackupunitCmd() *core.Command {
	backupUnitCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backupunit",
			Aliases:          []string{"b", "backup"},
			Short:            "BackupUnit Operations",
			Long:             "The sub-commands of `ionosctl compute backupunit` allow you to list, get, create, update, delete BackupUnits under your account.",
			TraverseChildren: true,
		},
	}
	backupUnitCmd.AddColsFlag(allCols)

	backupUnitCmd.AddCommand(BackupUnitListCmd())
	backupUnitCmd.AddCommand(BackupUnitGetCmd())
	backupUnitCmd.AddCommand(BackupUnitGetSsoUrlCmd())
	backupUnitCmd.AddCommand(BackupUnitCreateCmd())
	backupUnitCmd.AddCommand(BackupUnitUpdateCmd())
	backupUnitCmd.AddCommand(BackupUnitDeleteCmd())

	return core.WithConfigOverride(backupUnitCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
