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
		LongDesc:   "List all Target Groups in your contract. Target Groups are contract-wide (not scoped to a datacenter or ALB), so this returns every group regardless of which forwarding rules reference it. Use `ionosctl compute targetgroup target list` to see the backend servers inside a specific group.",
		Example:    `ionosctl compute targetgroup list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunTargetGroupList,
		InitClient: true,
	})

	return cmd
}
