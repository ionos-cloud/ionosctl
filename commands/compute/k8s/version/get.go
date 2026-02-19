package version

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func K8sVersionGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get Kubernetes Default Version",
		LongDesc:   "Use this command to retrieve the current default Kubernetes version for Clusters and NodePools.",
		Example:    "ionosctl compute k8s version get",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sVersionGet,
		InitClient: true,
	})

	return cmd
}
