package version

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
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
		Example:    "ionosctl k8s version get",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunK8sVersionGet,
		InitClient: true,
	})

	return cmd
}
