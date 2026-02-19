package snapshot

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func SnapshotListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "snapshot",
		Resource:   "snapshot",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Snapshots",
		LongDesc:   "Use this command to get a list of Snapshots.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.SnapshotsFiltersUsage(),
		Example:    "ionosctl compute snapshot list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunSnapshotList,
		InitClient: true,
	})

	return cmd
}
