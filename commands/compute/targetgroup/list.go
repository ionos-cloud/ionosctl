package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func TargetGroupListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "targetgroup",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Target Groups",
		LongDesc:   "Use this command to get a list of Target Groups.",
		Example:    `ionosctl targetgroup list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTargetGroupList,
		InitClient: true,
	})

	return cmd
}
