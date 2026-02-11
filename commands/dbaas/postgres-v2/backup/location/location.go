package location

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

var (
	defaultBackupLocationCols = []string{"LocationId", "Location"}
	allBackupLocationCols     = []string{"LocationId", "Location"}
)

func BackupLocationCmd() *core.Command {
	locationCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "location",
			Aliases:          []string{"loc"},
			Short:            "PostgreSQL Backup Location Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres backup location` allow you to list and get DBaaS PostgreSQL Backup Locations.",
			TraverseChildren: true,
		},
	}

	locationCmd.AddCommand(BackupLocationListCmd())
	locationCmd.AddCommand(BackupLocationGetCmd())

	return locationCmd
}
