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
		LongDesc:   "List the Managed Kubernetes clusters (control planes) in your contract, showing each cluster's ID, name, Kubernetes version, state, maintenance window and public/private flag.\n\nYou can filter the results using the `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.K8sClustersFiltersUsage(),
		Example:    "ionosctl compute k8s cluster list",
		PreCmdRun:  PreRunK8sClusterList,
		CmdRun:     RunK8sClusterList,
		InitClient: true,
	})

	return cmd
}
