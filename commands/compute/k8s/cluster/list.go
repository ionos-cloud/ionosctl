package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func K8sClusterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes Clusters",
		LongDesc:   "Use this command to get a list of existing Kubernetes Clusters.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.K8sClustersFiltersUsage(),
		Example:    "ionosctl k8s cluster list",
		PreCmdRun:  PreRunK8sClusterList,
		CmdRun:     RunK8sClusterList,
		InitClient: true,
	})

	return cmd
}
