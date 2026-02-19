package version

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func K8sVersionListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "version",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes Versions",
		LongDesc:   "Use this command to retrieve all available Kubernetes versions.",
		Example:    "ionosctl compute k8s version list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sVersionList,
		InitClient: true,
	})

	return cmd
}
