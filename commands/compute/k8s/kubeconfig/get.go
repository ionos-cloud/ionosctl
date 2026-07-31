package kubeconfig

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func K8sKubeconfigGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "kubeconfig",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get the kubeconfig file for a Kubernetes Cluster",
		LongDesc:  "Retrieve the kubeconfig file for a Managed Kubernetes cluster and print it to stdout. It carries the API-server endpoint and credentials that kubectl needs. Save it to a file and point kubectl at it via the KUBECONFIG environment variable (see the examples).\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example: `# Save the kubeconfig to a file and use it with kubectl
ionosctl compute k8s kubeconfig get --cluster-id CLUSTER_ID > kubeconfig.yaml
export KUBECONFIG=$PWD/kubeconfig.yaml
kubectl get nodes`,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sKubeconfigGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(constants.FlagClusterId, cloudapiv6.ArgIdShort, "", cloudapiv6.K8sClusterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
