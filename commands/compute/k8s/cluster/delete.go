package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func K8sClusterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Kubernetes Cluster",
		LongDesc: `Delete a Managed Kubernetes cluster (control plane).

The cluster must contain no node pools before it can be deleted - delete every
node pool first (see ` + "`ionosctl compute k8s nodepool delete`" + `), otherwise the
request is rejected. Deleting the cluster does not delete the Data Centers its
worker Nodes were placed in.

Use ` + "`" + `--wait` + "`" + ` (` + "`" + `-w` + "`" + `) to block until the deletion completes.

Required values to run command:

* K8s Cluster Id`,
		Example:    "ionosctl compute k8s cluster delete --cluster-id CLUSTER_ID",
		PreCmdRun:  PreRunK8sClusterDelete,
		CmdRun:     RunK8sClusterDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(constants.FlagClusterId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every Kubernetes cluster in the contract (each must already be empty of node pools)")

	return cmd
}
