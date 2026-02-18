package resource

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ResourceListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "resource",
		Resource:   "resource",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Resources",
		LongDesc:   "Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.",
		Example:    `ionosctl resource list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunResourceList,
		InitClient: true,
	})

	return cmd
}
