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
			Use:     "backupunit",
			Aliases: []string{"b", "backup"},
			Short:   "BackupUnit Operations",
			Long: `A BackupUnit is a named storage container plus a login used by the IONOS Managed Backup (MBaaS) agent to store server backups.

Domain model:
  * A BackupUnit belongs to your contract. Its "name" is GLOBALLY UNIQUE across all IONOS contracts and, combined with your contract number, forms the backup login: CONTRACT_NUMBER-BACKUP_UNIT_NAME. The name CANNOT be changed after creation.
  * The "password" set at creation is the login secret used to register the backup agent. It is write-only: the IONOS CLOUD API never returns it, so record it when you create the unit (only --password and --email can be updated later).
  * The "email" receives service reports from the backup system; it is independent of your IONOS CLOUD API username.

Backups themselves (backup plans, schedules, retention) are managed inside the backup web console, not through this API. Use ` + "`backupunit get-sso-url`" + ` to obtain a single-sign-on link into that console (https://backup.ionos.com).

The sub-commands of ` + "`ionosctl compute backupunit`" + ` let you list, get, create, update and delete BackupUnits, and fetch the console SSO URL.`,
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
