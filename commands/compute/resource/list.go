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
		ShortDesc:  "List all shareable resources on the contract",
		LongDesc:   "List every resource on the contract that can be shared with a Group (datacenters, snapshots, images, IP blocks, PCCs, backup units, Kubernetes clusters), across all types. To narrow the list to a single type, use `ionosctl compute resource get --resource-type <type>`.",
		Example:    `ionosctl compute resource list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunResourceList,
		InitClient: true,
	})

	return cmd
}
