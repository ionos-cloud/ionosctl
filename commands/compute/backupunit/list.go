package backupunit

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func BackupUnitListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "backupunit",
		Resource:   "backupunit",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List BackupUnits",
		LongDesc:   "Use this command to get a list of existing BackupUnits available on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.BackupUnitsFiltersUsage(),
		Example:    `ionosctl backupunit list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunBackupUnitList,
		InitClient: true,
	})

	return cmd
}
