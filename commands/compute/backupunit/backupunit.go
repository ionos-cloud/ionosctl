package backupunit

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const backupUnitNote = "NOTE: To login with backup agent use: https://backup.ionos.com, with CONTRACT_NUMBER-BACKUP_UNIT_NAME and BACKUP_UNIT_PASSWORD!"

var (
	defaultBackupUnitCols   = []string{"BackupUnitId", "Name", "Email", "State"}
	defaultBackupUnitSSOUrl = []string{"BackupUnitSsoUrl"}
)

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
	globalFlags := backupUnitCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultBackupUnitCols, tabheaders.ColsMessage(defaultBackupUnitCols))
	_ = viper.BindPFlag(core.GetFlagName(backupUnitCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = backupUnitCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultBackupUnitCols, cobra.ShellCompDirectiveNoFileComp
	})

	backupUnitCmd.AddCommand(BackupUnitListCmd())
	backupUnitCmd.AddCommand(BackupUnitGetCmd())
	backupUnitCmd.AddCommand(BackupUnitGetSsoUrlCmd())
	backupUnitCmd.AddCommand(BackupUnitCreateCmd())
	backupUnitCmd.AddCommand(BackupUnitUpdateCmd())
	backupUnitCmd.AddCommand(BackupUnitDeleteCmd())

	return core.WithConfigOverride(backupUnitCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
